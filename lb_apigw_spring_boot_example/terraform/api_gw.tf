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

resource "google_service_account" "api_gw_sa" {
  project      = var.project_id
  account_id   = "helloworld-api-gateway"
  display_name = "API Gateway Service Account"

  depends_on = [
    google_project_service.project
  ]
}

resource "google_api_gateway_api" "api" {
  provider = google-beta
  project  = var.project_id
  api_id   = "helloworld-api"

  depends_on = [
    google_project_service.project
  ]
}

resource "google_api_gateway_api_config" "api_cfg" {
  provider = google-beta
  project  = var.project_id
  api      = google_api_gateway_api.api.api_id

  api_config_id = "helloworld-api-cfg-4"

  openapi_documents {
    document {
      path = "spec.yaml"
      contents = textencodebase64(templatefile("api-spec.yaml", {
        "BACKEND_SERVICE_NAME" = google_cloud_run_service.spring-app.status[0].url
        "OAUTH_AUTHORIZATION_URL" = var.oauth_authorization_url
        "OAUTH_ISSUER" = var.oauth_client_issuer
        "OAUTH_CLIENT_ID" = var.oauth_client_id
        "OAUTH_CLIENT_SECRET" = var.oauth_client_secret
        "OAUTH_AUDIENCE" = var.oauth_audience
        "OAUTH_CLIENT_JWKS" = var.oauth_client_jwks
      }), "utf-8")
    }
  }

  gateway_config {
    backend_config {
      google_service_account = google_service_account.api_gw_sa.email
    }
  }
  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    google_api_gateway_api.api,
    google_project_service.project
  ]
}

resource "google_api_gateway_gateway" "api_gw" {
  provider   = google-beta
  project    = var.project_id
  region     = var.region
  api_config = google_api_gateway_api_config.api_cfg.id
  gateway_id = "api-gw"
  depends_on = [
    google_api_gateway_api_config.api_cfg,
    google_api_gateway_api.api,
    google_project_service.project
  ]
}

resource "google_compute_region_network_endpoint_group" "api_gw_neg" {
  provider              = google-beta
  name                  = "api-gw-neg"
  network_endpoint_type = "SERVERLESS"
  project               = var.project_id
  region                = var.region
  serverless_deployment {
    platform = "apigateway.googleapis.com"
    resource = google_api_gateway_gateway.api_gw.gateway_id
  }
}

module "api-gw-lb-http" {
  source  = "GoogleCloudPlatform/lb-http/google//modules/serverless_negs"
  version = "~> 6.2"
  name    = "api-helloworld"
  project = var.project_id

  ssl                             = true
  managed_ssl_certificate_domains = [var.api_domain]
  https_redirect                  = true
  backends = {
    default = {
      description = null
      groups = [
        {
          group = google_compute_region_network_endpoint_group.api_gw_neg.self_link
        }
      ]
      enable_cdn              = false
      security_policy         = null
      custom_request_headers  = null
      custom_response_headers = null

      iap_config = {
        enable               = false
        oauth2_client_id     = ""
        oauth2_client_secret = ""
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

output "api_public_ip" {
  value = module.api-gw-lb-http.external_ip
}