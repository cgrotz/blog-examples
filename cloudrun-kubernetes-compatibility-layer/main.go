// Copyright 2022 Google LLC
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
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"

	"github.com/gofiber/fiber/v2"
)

func Api(name string, version string) *fiber.Map {

	groupVersion := fmt.Sprintf("%s/%s", name, version)
	versionMap := &fiber.Map{
		"groupVersion": groupVersion,
		"version":      version,
	}

	return &fiber.Map{
		"name":     name,
		"versions": []*fiber.Map{versionMap},
		"preferredVersion": &fiber.Map{
			"groupVersion": groupVersion,
			"version":      version,
		},
	}
}

func main() {
	var err error
	app := fiber.New()
	app.Use(logger.New())
	app.Use(func(c *fiber.Ctx) error {
		auth := string(c.Request().Header.Peek("Authorization"))
		if auth == "" {
			c.Response().Header.Add("WWW-Authenticate", "Basic realm=\"Default Realm\"")
			return c.SendStatus(401)
		} else {
			return c.Next()
		}
	})
	app.Get("/api", func(c *fiber.Ctx) error {
		requestContentType := string(c.Request().Header.ContentType())
		log.Printf("Got API request content-type=%s \n", requestContentType)
		return c.Status(200).JSON(
			&fiber.Map{
				"kind":     "APIVersions",
				"versions": []string{"v1"},
			},
		)
	})
	app.Get("/version", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(
			&fiber.Map{
				"major": "1",
				"minor": "21",
			},
		)
	})
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(
			&fiber.Map{
				"kind":       "APIResourceList",
				"apiVersion": "v1",
				"resources":  []*fiber.Map{},
			},
		)
	})

	app.Get("/apis", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(
			&fiber.Map{
				"kind":       "APIGroupList",
				"apiVersion": "v1",
				"groups": []*fiber.Map{
					Api("serving.knative.dev", "v1"),
				},
			},
		)
	})

	app.Get("/apis/serving.knative.dev/v1", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(
			&fiber.Map{
				"kind":         "APIResourceList",
				"apiVersion":   "v1",
				"groupVersion": "serving.knative.dev/v1",
				"resources": []*fiber.Map{
					&fiber.Map{
						"name":         "services",
						"singularName": "service",
						"namespaced":   true,
						"kind":         "Service",
						"verbs": []string{
							"delete",
							"deletecollection",
							"get",
							"list",
							"patch",
							"create",
							"update",
							"watch",
						},
						"shortNames": []string{
							"kservice",
							"ksvc",
						},
						"categories": []string{
							"all",
							"knative",
							"serving",
						},
					},
					&fiber.Map{
						"name":         "services/status",
						"singularName": "",
						"namespaced":   true,
						"kind":         "Service",
						"verbs": []string{
							"get",
							"patch",
							"update",
						},
					},
					&fiber.Map{
						"name":         "routes",
						"singularName": "route",
						"namespaced":   true,
						"kind":         "Route",
						"verbs": []string{
							"delete",
							"deletecollection",
							"get",
							"list",
							"patch",
							"create",
							"update",
							"watch",
						},
						"shortNames": []string{
							"rt",
						},
						"categories": []string{
							"all",
							"knative",
							"serving",
						},
					},
					&fiber.Map{
						"name":         "routes/status",
						"singularName": "",
						"namespaced":   true,
						"kind":         "Route",
						"verbs": []string{
							"get",
							"patch",
							"update",
						},
					},
					&fiber.Map{
						"name":         "revisions",
						"singularName": "revision",
						"namespaced":   true,
						"kind":         "Revision",
						"verbs": []string{
							"delete",
							"deletecollection",
							"get",
							"list",
							"patch",
							"create",
							"update",
							"watch",
						},
						"shortNames": []string{
							"rev",
						},
						"categories": []string{
							"all",
							"knative",
							"serving",
						},
					},
					&fiber.Map{
						"name":         "revisions/status",
						"singularName": "",
						"namespaced":   true,
						"verbs": []string{
							"get",
							"patch",
							"update",
						},
					},
					&fiber.Map{
						"name":         "configurations",
						"singularName": "configuration",
						"namespaced":   true,
						"kind":         "Configuration",
						"verbs": []string{
							"delete",
							"deletecollection",
							"get",
							"list",
							"patch",
							"create",
							"update",
							"watch",
						},
						"shortNames": []string{
							"config",
							"cfg",
						},
						"categories": []string{
							"all",
							"knative",
							"serving",
						},
					},
					&fiber.Map{
						"name":         "configurations/status",
						"singularName": "",
						"namespaced":   true,
						"kind":         "Configuration",
						"verbs": []string{
							"get",
							"patch",
							"update",
						},
					},
				},
			},
		)
	})

	app.All("/apis/serving.knative.dev/v1/namespaces/*", func(c *fiber.Ctx) error {
		url := apiHost() + c.Path()
		if err := proxy.Do(c, url); err != nil {
			return err
		}
		return nil
	})

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

func apiHost() string {
	if os.Getenv("CLOUD_RUN_API_HOST") != "" {
		return os.Getenv("CLOUD_RUN_API_HOST")
	} else {
		return "https://us-central1-run.googleapis.com:443"
	}
}
