// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/goombaio/namegenerator"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/cgrotz/minimal-sql-app/simple/model"
	"github.com/cgrotz/minimal-sql-app/simple/table"
	. "github.com/go-jet/jet/v2/mysql"
)

var db *sql.DB

func initTracer() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		log.Fatalf("texporter.NewExporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
		sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)

}

func main() {
	var err error
	initTracer()
	cfg := mysql.Config{
		User:                 os.Getenv("DATABASE_USER"),
		Passwd:               os.Getenv("DATABASE_PASSWORD"),
		Net:                  os.Getenv("DATABASE_NET"),
		Addr:                 os.Getenv("DATABASE_HOST"),
		DBName:               os.Getenv("DATABASE_NAME"),
		MultiStatements:      true,
		AllowNativePasswords: true,
	}

	driverName, err := otelsql.Register("mysql", semconv.DBSystemMySQL.Value.AsString())
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatalf("Failed to register the ocsql driver: %v", err)
	}
	db, err = sql.Open(driverName, cfg.FormatDSN())

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(5)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = MigrateDatabase(os.Getenv("DATABASE_NAME"), db)
	if err != nil {
		log.Fatal("Unable to initalize database")
	}
	authorCount := GetAuthorCount()
	if authorCount < 1 {
		log.Printf("%d authors in db yet, creating some more", authorCount)
		GenerateAuthors(db)
	} else {
		log.Printf("%d authors in db", authorCount)
	}
	app := fiber.New()

	app.Use(TracingMiddleware)

	// No static assets right now app.Static("/", "./public")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("IPAM Autopilot up and running 👋!")
	})

	app.Get("/authors", GetAuthors)

	var port int64
	if os.Getenv("PORT") != "" {
		port, err = strconv.ParseInt(os.Getenv("PORT"), 10, 64)
		if err != nil {
			log.Panicf("Can't parse value of PORT env variable %s %v", os.Getenv("PORT"), err)
		}
	} else {
		port = 8080
	}
	app.Listen(fmt.Sprintf(":%d", port))
}

func GetAuthors(c *fiber.Ctx) error {
	tracer := otel.GetTracerProvider().Tracer("")
	ctx, span := tracer.Start(c.UserContext(), "authors.get")
	limitString := c.Query("limit")
	defer span.End()
	var limit int64
	var err error
	if limitString == "" {
		limit = -1
	} else {
		limit, err = strconv.ParseInt(limitString, 10, 64)
		if err != nil {
			log.Printf("Failed parsing limit query parameter '%s' error: %v", limitString, err)
			limit = -1
		}
	}
	authors, err := GetAuthorsFromDB(ctx, limit)
	if err != nil {
		return c.Status(503).JSON(&fiber.Map{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
	}
	//.JSON(results)
	/*
		var results []*fiber.Map
		return c.Status(200).JSON(&fiber.Map{
			"success": true,
			"message": "wow",
		})*/
	//return c.Status(200).JSON(SendString("IPAM Autopilot up and running 👋!")
	return json.NewEncoder(c.Response().BodyWriter()).Encode(authors)
}

func GenerateAuthors(db *sql.DB) {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	for i := 0; i < 10000; i++ {
		name := nameGenerator.Generate()
		CreateAuthorsInDb(name)
	}
}

func CreateAuthorsInDb(name string) (int64, error) {
	stmt := table.Authors.INSERT(table.Authors.Name).VALUES(name)

	res, err := stmt.Exec(db)
	if err != nil {
		log.Fatal(err)
	}
	author_id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return author_id, nil
}

func GetAuthorsFromDB(ctx context.Context, limit int64) ([]model.Authors, error) {
	tracer := otel.GetTracerProvider().Tracer("")
	_, span := tracer.Start(ctx, "authors.db.get")
	defer span.End()
	stmt := SELECT(
		table.Authors.AuthorID.AS("authors.author_id"),
		table.Authors.Name.AS("authors.name"),
	).FROM(
		table.Authors,
	)
	if limit > -1 {
		stmt = stmt.LIMIT(1000)
	}

	var dest []model.Authors
	err := stmt.QueryContext(ctx, db, &dest)
	if err != nil {
		log.Fatal(err)
	}
	return dest, nil
}

type AuthorCount struct {
	Count int64
}

func GetAuthorCount() int64 {
	stmt := SELECT(
		COUNT(STAR).AS("author_count.count"),
	).FROM(
		table.Authors,
	)
	var dest AuthorCount
	err := stmt.Query(db, &dest)
	if err != nil {
		log.Fatal(err)
	}
	return dest.Count
}

func TracingMiddleware(c *fiber.Ctx) error {
	// "X-Cloud-Trace-Context: TRACE_ID/SPAN_ID;o=TRACE_TRUE"
	traceHeader := string(c.Request().Header.Peek("X-Cloud-Trace-Context"))
	if traceHeader != "" {
		traceIdString := strings.Split(traceHeader, "/")[0]
		spanIdString := strings.Split(strings.Split(traceHeader, "/")[1], ";")[0]
		sampling := strings.Index(traceHeader, ";o=1")

		traceId, err := oteltrace.TraceIDFromHex(traceIdString)
		if err != nil {
			fmt.Printf("Unable to extract trace from header %s %v", traceHeader, err)
		}
		spanId, err := oteltrace.SpanIDFromHex(spanIdString)
		if err != nil {
			fmt.Printf("Unable to extract span from header %s %v", traceHeader, err)
		}
		log.Printf("TraceContext present header=%s, traceId=%s, spanId=%s\n", traceHeader, traceIdString, spanIdString)
		spanContext := oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: traceId,
			SpanID:  spanId,
			Remote:  sampling == -1,
		})
		ctx := oteltrace.ContextWithSpanContext(c.Context(), spanContext)

		tracer := otel.GetTracerProvider().Tracer("")
		ctx, span := tracer.Start(ctx, fmt.Sprintf("Fiber %s %s", c.Method(), c.Path()))
		defer span.End()
		c.SetUserContext(ctx)
		return c.Next()
	} else {
		log.Printf("No TraceContext present\n")
		return c.Next()
	}
}
