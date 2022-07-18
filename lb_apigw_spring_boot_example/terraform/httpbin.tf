/**
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */



resource "google_service_account" "httpbin" {
  project      = var.project_id
  account_id   = "httpbin-helloworld"
  display_name = "httpbin Service Account"

  depends_on = [
    google_project_service.project
  ]
}

resource "google_cloud_run_service" "httpbin" {
  name     = "helloworld-httpbin"
  location = var.region
  project  = var.project_id

  template {
    spec {
      containers {
        image = "mirror.gcr.io/kennethreitz/httpbin"
        ports {
          container_port = 80
        }
      }
      service_account_name = google_service_account.httpbin.email
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
  metadata {
    annotations = {
      "run.googleapis.com/ingress" = "all" // Is needed for now
    }
  }
}

resource "google_cloud_run_service_iam_member" "member" {
  location = google_cloud_run_service.httpbin.location
  project  = google_cloud_run_service.httpbin.project
  service  = google_cloud_run_service.httpbin.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.spring-app.email}"
}