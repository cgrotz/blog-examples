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

resource "google_identity_platform_oauth_idp_config" "oauth_idp_config" {
  name          = "oidc.oauth-idp-config"
  display_name  = "OAuth provider"
  project       = var.project_id
  enabled       = true
  client_id     = var.oauth_client_id
  client_secret = var.oauth_client_secret
  issuer        = var.oauth_client_issuer

  depends_on = [
    google_project_service.project
  ]
}

/*
When https://github.com/GoogleCloudPlatform/magic-modules/pull/6249 is merged and released

resource "google_apikeys_key" "loginUi" {
  name         = "loginUi"
  display_name = "loginUi"
  project      = var.project_id
  restrictions {
    browser_key_restrictions {
      allowed_referrers = ["https://iap.googleapis.com/v1/oauth/clientIds/${google_iap_client.project_client.client_id}:handleRedirect"]
    }
  }
}

resource "google_iap_settings" "settings" {
  backend_service_id = module.ui-lb-http.o
  access_settings {
    gcip_settings {
      tenant_ids = [ "_${data.google_project.project.number}" ]
    }
    loginPageUri = "https://${var.ui_domain}/login/index.html?apiKey=${google_apikeys_key.loginUi.key_string}"
  }
  depends_on = [
    google_project_service.project
  ]
}*/