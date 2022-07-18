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

resource "google_service_account" "spring-app" {
  project      = var.project_id
  account_id   = "spring-app"
  display_name = "Spring App Service Account"

  depends_on = [
    google_project_service.project
  ]
}

resource "google_cloud_run_service" "spring-app" {
  name     = "spring-app"
  location = var.region
  project  = var.project_id

  template {
    spec {
      containers {
        image = var.spring_app_image_name
        ports {
          container_port = 8080
        }
        env {
          name  = "BACKEND_SERVICE_NAME"
          value = google_cloud_run_service.httpbin.status[0].url
        }
        env {
          name  = "OIDC_ISSUER"
          value = var.oauth_client_issuer
        }
        env {
          name  = "OIDC_JWKS"
          value = var.oauth_client_jwks
        }
      }
      service_account_name = google_service_account.spring-app.email
    }
  }

  autogenerate_revision_name = true
  traffic {
    percent         = 100
    latest_revision = true
  }
  metadata {
    annotations = {
      "run.googleapis.com/ingress"      = "all" // Is needed for now
      "run.googleapis.com/launch-stage" = "BETA"
    }
  }
}

resource "google_cloud_run_service_iam_member" "spring-app-member" {
  location = google_cloud_run_service.spring-app.location
  project  = google_cloud_run_service.spring-app.project
  service  = google_cloud_run_service.spring-app.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.api_gw_sa.email}"
}

resource "google_cloud_run_service" "spring-app-internal" {
  name     = "spring-app-internal"
  location = var.region
  project  = var.project_id

  template {
    spec {
      containers {
        image = var.spring_app_image_name
        ports {
          container_port = 8080
        }
        env {
          name  = "BACKEND_SERVICE_NAME"
          value = google_cloud_run_service.httpbin.status[0].url
        }
        env {
          name  = "OIDC_ISSUER"
          value = var.oauth_client_issuer
        }
        env {
          name  = "OIDC_JWKS"
          value = var.oauth_client_jwks
        }
      }
      service_account_name = google_service_account.spring-app.email
    }
  }
  autogenerate_revision_name = true
  traffic {
    percent         = 100
    latest_revision = true
  }
  metadata {
    annotations = {
      "run.googleapis.com/ingress"      = "internal-and-cloud-load-balancing" // Is needed for now
      "run.googleapis.com/launch-stage" = "BETA"
    }
  }
}

resource "google_cloud_run_service_iam_member" "spring-app-internal-member" {
  location = google_cloud_run_service.spring-app-internal.location
  project  = google_cloud_run_service.spring-app-internal.project
  service  = google_cloud_run_service.spring-app-internal.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}