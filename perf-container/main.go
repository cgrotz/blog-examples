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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var preRequestDelay = flag.Int("preRequestDelay", 0, "Delay in Milliseconds before the request is passed to the backend")
var postRequestDelay = flag.Int("postRequestDelay", 0, "Delay in Milliseconds after the request is passed to the backend")
var processingTime = flag.Int("processingTime", 0, "Time in Milliseconds before for fake processing before the request is responded to")
var startupDelay = flag.Int("startupDelay", 0, "Delay in Milliseconds before the container starts up")
var reverseProxyDestination = flag.String("remote", "", "Backend for the reverse proxy")
var runAsReverseProxy = flag.Bool("proxy", false, "Should the app run in reverse proxy mode")
var port = flag.Int("port", 8080, "Server port for the app")
var tracing = flag.Bool("tracing", true, "Tracing enabled; defaults to true")
var explicitError = flag.Bool("error", false, "Explicitly throw an error before starting the HTTP server; defaults to false")
var project string

func main() {
	flag.Parse()
	tracing_temp := GetBoolValue("TRACING", *tracing)
	tracing = &tracing_temp
	if *tracing == true {
		if err := initTracer(); err != nil {
			log.Fatalf("Failed creating tracer %v", err)
		}
	}

	// Get project ID from metadata server
	project = ""
	if *tracing {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
		req.Header.Set("Metadata-Flavor", "Google")
		res, err := client.Do(req)
		if err == nil {
			defer res.Body.Close()
			if res.StatusCode == 200 {
				responseBody, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatal(err)
				}
				project = string(responseBody)
			}
		}
	}

	structuredLogging(map[string]interface{}{
		"type":  "event",
		"event": "STARTING",
	})

	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		structuredLogging(map[string]interface{}{
			"type":  "event",
			"event": "STOPPING",
		})
		close(done)
		os.Exit(0)
	}()

	preRequestDelay = GetIntValueFromEnvOrUseFlag("PRE_REQUEST_DELAY", *preRequestDelay)
	postRequestDelay = GetIntValueFromEnvOrUseFlag("POST_REQUEST_DELAY", *postRequestDelay)
	processingTime = GetIntValueFromEnvOrUseFlag("PROCESSING_TIME", *processingTime)
	explicitError = GetBoolValueFromEnvOrUseFlag("EXPLICIT_ERROR", *explicitError)

	time.Sleep(time.Duration(GetIntValueFromEnv("STARTUP_DELAY", *startupDelay) * int(time.Millisecond)))

	if GetBoolValue("RUN_AS_REVERSE_PROXY", *runAsReverseProxy) {
		destination := GetStringValue("REVERSE_PROXY_DESTINATION", *reverseProxyDestination)
		if destination == "" {
			log.Panicln("In proxy mode destination needs to be set")
		}
		remote, err := url.Parse(destination)
		if err != nil {
			log.Panicf("destination is not a parsable URL %v", err)
		}
		handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				localPreRequestDelay, err := GetQueryOrDefault(r, "pre_request_delay", *preRequestDelay)
				if err != nil {
					fmt.Fprintf(w, "%v", err)
					return
				}
				localPostRequestDelay, err := GetQueryOrDefault(r, "post_request_delay", *postRequestDelay)
				if err != nil {
					fmt.Fprintf(w, "%v", err)
					return
				}

				if *tracing {
					tracer := otel.GetTracerProvider().Tracer("")
					_, spanOuter := tracer.Start(r.Context(), "proxying-outer")
					structuredLogging(map[string]interface{}{
						"logging.googleapis.com/trace":  fmt.Sprintf("projects/%s/traces/%s", project, spanOuter.SpanContext().TraceID().String()),
						"logging.googleapis.com/spanId": spanOuter.SpanContext().SpanID().String(),
						"time":                          time.Now().UnixNano() / int64(time.Millisecond),
						"event":                         "proxy-start",
					})
					defer spanOuter.End()
				} else {
					structuredLogging(map[string]interface{}{
						"time":  time.Now().UnixNano() / int64(time.Millisecond),
						"event": "proxy-start",
					})
				}

				time.Sleep(time.Duration(localPreRequestDelay * int(time.Millisecond)))
				r.Host = remote.Host
				if *tracing {
					tracer := otel.GetTracerProvider().Tracer("")
					_, spanInner := tracer.Start(r.Context(), "proxying-inner")
					start := time.Now()
					structuredLogging(map[string]interface{}{
						"logging.googleapis.com/trace":  fmt.Sprintf("projects/%s/traces/%s", project, spanInner.SpanContext().TraceID().String()),
						"logging.googleapis.com/spanId": spanInner.SpanContext().SpanID().String(),
						"time":                          start.UnixNano() / int64(time.Millisecond),
						"event":                         "proxy-send",
					})
					p.ServeHTTP(w, r)
					structuredLogging(map[string]interface{}{
						"logging.googleapis.com/trace":  fmt.Sprintf("projects/%s/traces/%s", project, spanInner.SpanContext().TraceID().String()),
						"logging.googleapis.com/spanId": spanInner.SpanContext().SpanID().String(),
						"proxy_time":                    time.Since(start).Seconds(),
					})
					spanInner.End()
				} else {
					start := time.Now()
					p.ServeHTTP(w, r)
					structuredLogging(map[string]interface{}{
						"proxy_time": time.Since(start).Seconds(),
					})
				}
				time.Sleep(time.Duration(localPostRequestDelay * int(time.Millisecond)))
			}
		}
		if *tracing {
			otelHandler := otelhttp.NewHandler(http.HandlerFunc(handler(httputil.NewSingleHostReverseProxy(remote))), "proxy")
			http.Handle("/", otelHandler)
		} else {
			http.HandleFunc("/", handler(httputil.NewSingleHostReverseProxy(remote)))
		}
	} else {
		if *tracing {
			http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(LoopBack), "loopback"))
		} else {
			http.HandleFunc("/", http.HandlerFunc(LoopBack))
		}
	}
	if *explicitError {
		log.Fatal("Intentionally Stop")
	}
	structuredLogging(map[string]interface{}{
		"time":  time.Now().UnixNano() / int64(time.Millisecond),
		"event": "HTTP_READY",
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", GetIntValueFromEnv("PORT", *port)), nil))
}

func LoopBack(w http.ResponseWriter, r *http.Request) {
	localProcessingTime, err := GetQueryOrDefault(r, "processing_time", *processingTime)
	if err != nil {
		fmt.Fprintf(w, "%v", w)
		return
	}
	if *tracing {
		tracer := otel.GetTracerProvider().Tracer("")
		_, spanInner := tracer.Start(r.Context(), "loopback-outer")
		structuredLogging(map[string]interface{}{
			"logging.googleapis.com/trace":  fmt.Sprintf("projects/%s/traces/%s", project, spanInner.SpanContext().TraceID().String()),
			"logging.googleapis.com/spanId": spanInner.SpanContext().SpanID().String(),
			"time":                          time.Now().UnixNano() / int64(time.Millisecond),
			"event":                         "loopback",
		})
		defer spanInner.End()
		time.Sleep(time.Duration(localProcessingTime * int(time.Millisecond)))

		_, spanOuter := tracer.Start(r.Context(), "loopback-inner")
		defer spanOuter.End()
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%q", dump)
	} else {
		time.Sleep(time.Duration(localProcessingTime * int(time.Millisecond)))

		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%q", dump)
	}
}

func GetIntValueFromEnv(env string, defaultValue int) int {
	if os.Getenv(env) != "" {
		value, err := strconv.ParseInt(os.Getenv(env), 10, 64)
		if err != nil {
			log.Panicf("Can't parse value of %s env variable %s %v", env, os.Getenv(env), err)
		}
		return int(value)
	} else {
		return defaultValue
	}
}

func GetIntValueFromEnvOrUseFlag(env string, defaultValue int) *int {
	if os.Getenv(env) != "" {
		value, err := strconv.ParseInt(os.Getenv(env), 10, 64)
		if err != nil {
			log.Panicf("Can't parse value of %s env variable %s %v", env, os.Getenv(env), err)
		}
		intVal := int(value)
		return &intVal
	} else {
		return &defaultValue
	}
}

func GetStringValue(env string, defaultValue string) string {
	if os.Getenv(env) != "" {
		return os.Getenv(env)
	} else {
		return defaultValue
	}
}

func GetBoolValue(env string, defaultValue bool) bool {
	if os.Getenv(env) != "" {
		if os.Getenv(env) == "true" {
			return true
		} else {
			return false
		}
	} else {
		return defaultValue
	}
}

func GetBoolValueFromEnvOrUseFlag(env string, defaultValue bool) *bool {
	if os.Getenv(env) != "" {
		value, err := strconv.ParseBool(os.Getenv(env))
		if err != nil {
			log.Panicf("Can't parse value of %s env variable %s %v", env, os.Getenv(env), err)
		}
		return &value
	} else {
		return &defaultValue
	}
}

func GetQueryOrDefault(r *http.Request, queryParameterName string, defaultValue int) (int, error) {
	queryParameter := r.URL.Query().Get(queryParameterName)
	if queryParameter == "" {
		return defaultValue, nil
	} else {
		value, err := strconv.ParseInt(queryParameter, 10, 64)
		if err != nil {
			return -1, err
		}
		return int(value), nil
	}
}

func GetQueryOrDefaultString(r *http.Request, queryParameterName string, defaultValue string) string {
	queryParameter := r.URL.Query().Get(queryParameterName)
	if queryParameter == "" {
		return defaultValue
	} else {
		return queryParameter
	}
}

func structuredLogging(values map[string]interface{}) {
	content, _ := json.Marshal(values)
	fmt.Printf("%s\n", string(content))
}

func initTracer() error {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		return fmt.Errorf("unable to create tracer %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
		sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}
