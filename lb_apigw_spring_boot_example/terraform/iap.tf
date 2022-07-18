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

locals {
  tf-cr-lb = "helloworld-iap"
}

resource "google_compute_region_network_endpoint_group" "cloud_run_neg" {
  name                  = "cloud-run-neg"
  network_endpoint_type = "SERVERLESS"
  region                = var.region
  project               = data.google_project.project.project_id
  cloud_run {
    service = google_cloud_run_service.spring-app-internal.name
  }
}

module "ui-lb-http" {
  source  = "GoogleCloudPlatform/lb-http/google//modules/serverless_negs"
  version = "~> 6.2"
  name    = local.tf-cr-lb
  project = data.google_project.project.project_id

  ssl                             = true
  managed_ssl_certificate_domains = [var.ui_domain]
  https_redirect                  = true

  // Use custom url map.
  url_map        = google_compute_url_map.url-map.self_link
  create_url_map = false

  backends = {
    default = {
      description = null
      groups = [
        {
          group = google_compute_region_network_endpoint_group.cloud_run_neg.id
        }
      ]
      enable_cdn              = false
      security_policy         = google_compute_security_policy.api-policy.id
      custom_request_headers  = null
      custom_response_headers = null

      iap_config = {
        enable               = true
        oauth2_client_id     = google_iap_client.project_client.client_id
        oauth2_client_secret = google_iap_client.project_client.secret
      }
      log_config = {
        enable      = false
        sample_rate = null
      }
    }
  }

  depends_on = [
    google_project_service.project
  ]
}

resource "google_compute_security_policy" "api-policy" {
  provider = google-beta
  name     = "api-policy"
  project  = data.google_project.project.project_id

  adaptive_protection_config {
    layer_7_ddos_defense_config {
      enable = true
    }
  }

  depends_on = [
    google_project_service.project
  ]
}

resource "google_iap_client" "project_client" {
  display_name = "LB Client"
  brand        = "projects/${data.google_project.project.number}/brands/${data.google_project.project.number}"

  depends_on = [
    google_project_service.project
  ]
}

output "iap_redirect_url" {
  value = "https://iap.googleapis.com/v1/oauth/clientIds/${google_iap_client.project_client.client_id}:handleRedirect"
}

resource "google_compute_url_map" "url-map" {
  // note that this is the name of the load balancer
  name            = local.tf-cr-lb
  default_service = module.ui-lb-http.backend_services["default"].self_link
  project         = var.project_id

  host_rule {
    hosts        = ["*"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = module.ui-lb-http.backend_services["default"].self_link

    path_rule {
      paths = [
        "/${var.login_path}",
        "/${var.login_path}/*"

      ]
      service = google_compute_backend_bucket.website.self_link
    }
  }

  depends_on = [
    google_project_service.project
  ]
}

resource "google_compute_backend_bucket" "website" {
  name        = "website"
  description = "Contains static resources for the auth UI"
  bucket_name = google_storage_bucket.website.name
  project     = var.project_id
  enable_cdn  = false

  depends_on = [
    google_project_service.project
  ]
}

output "iap_public_ip" {
  value = module.ui-lb-http.external_ip
}