// Copyright 2021 Google LLC
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

variable "region" {
  default = "europe-west1"
  description = "Google Cloud region in which the resources should be created"
}

variable "service_account_for_setup" {
  description = "The service account that is impersonated to deploy the resources"
}

variable "project_id" {
  description = "ID of the project to which the resources should be deployed"
}

variable "spring_app_image_name" {
  description = "Container image that should be deployed to Cloud Run"
}

variable "api_domain" {
  description = "API domain for the API Load Balancer"
}

variable "ui_domain" {
  description = "UI domain for the API Load Balancer"
}

variable "login_path" {
  description = "URL path for the login page"
  type        = string
  default     = "login"
}

variable "oauth_authorization_url" {
  description = "Authorization URL of the OIDC Identity Provider"
}

variable "oauth_client_id" {
  description = "Client ID for the external OIDC Identity Provider"
}

variable "oauth_client_secret" {
  description = "Client Secret for the external OIDC Identity Provider"
}

variable "oauth_client_issuer" {
  description = "Issuer of the external OIDC Identity Provider (iss claim in the JWT)"
}

variable "oauth_client_jwks" {
  description = "JWKS of the external OIDC Identity Provider, they keys used to sign the JWT must be contained"
}

variable "oauth_audience" {
  description = "Audience from the JWT"
}

variable "brand_support_email" {
  description = "Brand support email for the Google IAP Brand"
}