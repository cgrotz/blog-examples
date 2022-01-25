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
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var preRequestDelay = flag.Int("preRequestDelay", 0, "Delay in Milliseconds before the request is passed to the backend")
var postRequestDelay = flag.Int("postRequestDelay", 0, "Delay in Milliseconds after the request is passed to the backend")
var processingTime = flag.Int("processingTime", 0, "Time in Milliseconds before for fake processing before the request is responded to")
var startupDelay = flag.Int("startupDelay", 0, "Delay in Milliseconds before the container starts up")
var reverseProxyDestination = flag.String("remote", "", "Backend for the reverse proxy")
var runAsReverseProxy = flag.Bool("proxy", false, "Should the app run in reverse proxy mode")
var port = flag.Int("port", 8080, "Server port for the app")

func main() {
	flag.Parse()

	structuredLogging(map[string]interface{}{
		"type":  "event",
		"event": "STARTED",
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
	}()

	preRequestDelay = GetIntValueFromEnvOrUseFlag("PRE_REQUEST_DELAY", *preRequestDelay)
	postRequestDelay = GetIntValueFromEnvOrUseFlag("POST_REQUEST_DELAY", *postRequestDelay)
	processingTime = GetIntValueFromEnvOrUseFlag("PROCESSING_TIME", *processingTime)

	time.Sleep(time.Duration(GetIntValueFromEnv("STARTUP_DELAY", *startupDelay) * int(time.Millisecond)))

	if GetBoolValue("RUN_AS_REVERSE_PROXY", *runAsReverseProxy) {
		destination := GetStringValue("REVERSE_PROXY_DESTINATION", *reverseProxyDestination)
		if destination == "" {
			log.Panicln("In proxy mode destination needs to be set")
		}
		remote, err := url.Parse(destination)
		if err != nil {
			log.Panicln("destination is not a parsable URL", err)
		}
		handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(time.Duration(*preRequestDelay * int(time.Millisecond)))
				r.Host = remote.Host
				p.ServeHTTP(w, r)
				time.Sleep(time.Duration(*postRequestDelay * int(time.Millisecond)))
			}
		}
		http.HandleFunc("/", handler(httputil.NewSingleHostReverseProxy(remote)))
	} else {
		http.HandleFunc("/", LoopBack)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", GetIntValueFromEnv("PORT", *port)), nil))
}

func LoopBack(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(*processingTime * int(time.Millisecond)))

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%q", dump)
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

func structuredLogging(values map[string]interface{}) {
	content, _ := json.Marshal(values)
	fmt.Printf("%s\n", string(content))
}
