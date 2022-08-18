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
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

func main() {

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		var url string
		if os.Getenv("DB_URL") == "" {
			url = os.Args[1]
		} else {
			url = os.Getenv("DB_URL")
		}
		log.Printf("Connecting to database with %s\n", url)
		db, err := sql.Open("postgres", url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}

		rows, err := db.Query("SELECT Id, Username, Age FROM Users")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		} else {
			var resp []User
			for rows.Next() {
				var id int64
				var name string
				var age int64
				err = rows.Scan(&id, &name, &age)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "%v", err)
					return
				}
				resp = append(resp, User{
					Id:   id,
					Name: name,
					Age:  age,
				})
			}

			jsonResp, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResp)
			return
		}
	})

	log.Printf("Listening on :8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
