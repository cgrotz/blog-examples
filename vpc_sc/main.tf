/**
 * Copyright 2021 Google LLC
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

locals {
  apis = ["iam.googleapis.com", "compute.googleapis.com", "run.googleapis.com", "apigateway.googleapis.com", "servicemanagement.googleapis.com", "servicecontrol.googleapis.com", "compute.googleapis.com", "iap.googleapis.com"]
}

data "google_project" "project" {
  project_id = var.project_id
}

resource "google_project_service" "project" {
  for_each = toset(local.apis)
  project  = data.google_project.project.project_id
  service  = each.key
  disable_on_destroy = false
}

resource "null_resource" docker_image {

  provisioner "local-exec" {
    command = <<EOT
docker pull kennethreitz/httpbin:latest
docker tag kennethreitz/httpbin:latest gcr.io/${var.project_id}/httpbin:latest
docker push gcr.io/${var.project_id}/httpbin:latest
EOT
  }
}

resource "null_resource" proxy_docker_image {

  provisioner "local-exec" {
    command = <<EOT
    docker build -t gcr.io/${var.project_id}/proxy:latest ./proxy
    docker push gcr.io/${var.project_id}/proxy:latest
EOT
  }
}
//  docker build -t gcr.io/cloud-armor-test-333010/proxy:1 ./proxy

resource "google_cloud_run_service" "default" {
  name     = "cloudrun-srv"
  location = var.region
  project  = var.project_id

  metadata {
    annotations = {
      "run.googleapis.com/ingress" :  "internal-and-cloud-load-balancing"
    }
  }

  template {
    spec {
      containers {
        image ="gcr.io/${var.project_id}/httpbin:latest"
        ports {
          name = "http1"
          container_port = 80
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.project,
    null_resource.docker_image,
  ]
}

resource "google_cloud_run_service" "proxy-default" {
  name     = "proxy-srv"
  location = var.region
  project  = var.project_id

  metadata {
    annotations = {
      "run.googleapis.com/ingress" :  "internal-and-cloud-load-balancing"
    }
  }

  template {
    spec {
      containers {
        image ="gcr.io/${var.project_id}/proxy:1"
        ports {
          name = "http1"
          container_port = 80
        }
        env {
          name = "BACKEND_URL"
          value = google_cloud_run_service.default.status[0].url
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.project,
    null_resource.proxy_docker_image,
  ]


}