# Example for using IAP and API Gateway with Cloud Run services

## Setup steps:
* Enable Identity Platform in your GCP project https://console.cloud.google.com/customer-identity/onboarding
* Create a Service Account that has sufficient permissions to deploy the required resources (some of the resources can only be configured/deployed using a service account)
* Configure the OAuth Consent screen https://console.cloud.google.com/apis/credentials/consent
* Build the login-ui (`npm install` and `npm run bundle`)
* Activate the Container Registry `gcloud services enable containerregistry.googleapis.com`
* Build the app container (`mvn clean compile jib:build -Dimage=gcr.io/<project_id>/spring-app:1`)
* Update variables of Terraform deployment (ideally in a `terraform.tfvars` file)
* Deploy the setup `terraform apply` (you probably need to run `terraform init` first)
* Map the UI and API domain to your DNS provider using the outputted IPs
* Check that the Identity Platform handler is configure n your OIDC Provider(Allowed Callback URLs: https://<project_id>.firebaseapp.com/__/auth/handler)
* Add the webui domain to `Authorized Domains` in Identity Platform settings https://pantheon.corp.google.com/customer-identity/settings
* Switch the IAP Loadbalancers IAP settings to external Identities https://pantheon.corp.google.com/iam-admin/iap
* Retrieve the API Key, and add the API Key and update the authDomain in the `login-ui/src/script.ts`
* Rebuild the login-ui (`npm install` and `npm run bundle`)
* Rerun the Terraform automation


## Terraform Documentation
<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | 4.28.0 |
| <a name="provider_google-beta"></a> [google-beta](#provider\_google-beta) | 4.28.0 |
| <a name="provider_null"></a> [null](#provider\_null) | 3.1.1 |
| <a name="provider_random"></a> [random](#provider\_random) | 3.3.2 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_api-gw-lb-http"></a> [api-gw-lb-http](#module\_api-gw-lb-http) | GoogleCloudPlatform/lb-http/google//modules/serverless_negs | ~> 6.2 |
| <a name="module_ui-lb-http"></a> [ui-lb-http](#module\_ui-lb-http) | GoogleCloudPlatform/lb-http/google//modules/serverless_negs | ~> 6.2 |

## Resources

| Name | Type |
|------|------|
| [google-beta_google_api_gateway_api.api](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/google_api_gateway_api) | resource |
| [google-beta_google_api_gateway_api_config.api_cfg](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/google_api_gateway_api_config) | resource |
| [google-beta_google_api_gateway_gateway.api_gw](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/google_api_gateway_gateway) | resource |
| [google-beta_google_compute_region_network_endpoint_group.api_gw_neg](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/google_compute_region_network_endpoint_group) | resource |
| [google-beta_google_compute_security_policy.api-policy](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/google_compute_security_policy) | resource |
| [google_cloud_run_service.httpbin](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service) | resource |
| [google_cloud_run_service.spring-app](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service) | resource |
| [google_cloud_run_service.spring-app-internal](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service) | resource |
| [google_cloud_run_service_iam_member.member](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service_iam_member) | resource |
| [google_cloud_run_service_iam_member.spring-app-internal-member](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service_iam_member) | resource |
| [google_cloud_run_service_iam_member.spring-app-member](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service_iam_member) | resource |
| [google_compute_backend_bucket.website](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_backend_bucket) | resource |
| [google_compute_region_network_endpoint_group.cloud_run_neg](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_region_network_endpoint_group) | resource |
| [google_compute_url_map.url-map](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_url_map) | resource |
| [google_iap_client.project_client](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/iap_client) | resource |
| [google_identity_platform_oauth_idp_config.oauth_idp_config](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/identity_platform_oauth_idp_config) | resource |
| [google_project_service.project](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/project_service) | resource |
| [google_service_account.api_gw_sa](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/service_account) | resource |
| [google_service_account.httpbin](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/service_account) | resource |
| [google_service_account.spring-app](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/service_account) | resource |
| [google_storage_bucket.website](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket) | resource |
| [google_storage_bucket_iam_binding.binding](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket_iam_binding) | resource |
| [null_resource.upload_folder_content](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource) | resource |
| [random_id.random_id](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/id) | resource |
| [google_project.project](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/project) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_api_domain"></a> [api\_domain](#input\_api\_domain) | API domain for the API Load Balancer | `any` | n/a | yes |
| <a name="input_brand_support_email"></a> [brand\_support\_email](#input\_brand\_support\_email) | Brand support email for the Google IAP Brand | `any` | n/a | yes |
| <a name="input_login_path"></a> [login\_path](#input\_login\_path) | URL path for the login page | `string` | `"login"` | no |
| <a name="input_oauth_audience"></a> [oauth\_audience](#input\_oauth\_audience) | Audience from the JWT | `any` | n/a | yes |
| <a name="input_oauth_authorization_url"></a> [oauth\_authorization\_url](#input\_oauth\_authorization\_url) | Authorization URL of the OIDC Identity Provider | `any` | n/a | yes |
| <a name="input_oauth_client_id"></a> [oauth\_client\_id](#input\_oauth\_client\_id) | Client ID for the external OIDC Identity Provider | `any` | n/a | yes |
| <a name="input_oauth_client_issuer"></a> [oauth\_client\_issuer](#input\_oauth\_client\_issuer) | Issuer of the external OIDC Identity Provider (iss claim in the JWT) | `any` | n/a | yes |
| <a name="input_oauth_client_jwks"></a> [oauth\_client\_jwks](#input\_oauth\_client\_jwks) | JWKS of the external OIDC Identity Provider, they keys used to sign the JWT must be contained | `any` | n/a | yes |
| <a name="input_oauth_client_secret"></a> [oauth\_client\_secret](#input\_oauth\_client\_secret) | Client Secret for the external OIDC Identity Provider | `any` | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | ID of the project to which the resources should be deployed | `any` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | Google Cloud region in which the resources should be created | `string` | `"europe-west1"` | no |
| <a name="input_service_account_for_setup"></a> [service\_account\_for\_setup](#input\_service\_account\_for\_setup) | The service account that is impersonated to deploy the resources | `any` | n/a | yes |
| <a name="input_spring_app_image_name"></a> [spring\_app\_image\_name](#input\_spring\_app\_image\_name) | Container image that should be deployed to Cloud Run | `any` | n/a | yes |
| <a name="input_ui_domain"></a> [ui\_domain](#input\_ui\_domain) | UI domain for the API Load Balancer | `any` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_api_public_ip"></a> [api\_public\_ip](#output\_api\_public\_ip) | n/a |
| <a name="output_iap_public_ip"></a> [iap\_public\_ip](#output\_iap\_public\_ip) | n/a |
| <a name="output_iap_redirect_url"></a> [iap\_redirect\_url](#output\_iap\_redirect\_url) | n/a |
<!-- END_TF_DOCS -->
