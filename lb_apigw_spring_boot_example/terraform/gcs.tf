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

resource "random_id" "random_id" {
  byte_length = 8
}

resource "google_storage_bucket" "website" {
  name     = "hello-world-${random_id.random_id.hex}"
  location = "EU"
  project  = var.project_id

  website {
    main_page_suffix = "index.html"
    not_found_page   = "${var.login_path}/404.html"
  }

  // delete bucket and contents on destroy.
  force_destroy = true

  depends_on = [
    google_project_service.project
  ]
}

resource "google_storage_bucket_iam_binding" "binding" {
  bucket = google_storage_bucket.website.name
  role = "roles/storage.objectViewer"
  members = [
    "allUsers",
  ]
}

resource "null_resource" "upload_folder_content" {
  //only upload files if any change
  triggers = {
    file_hashes = jsonencode({
      for fn in fileset("../login-ui/public", "**") :
      fn => filesha256("../login-ui/public/${fn}")
    })
  }
  provisioner "local-exec" {
    command = "gsutil cp -r ../login-ui/public/* gs://${google_storage_bucket.website.name}/${var.login_path}"
  }
  depends_on = [
    google_storage_bucket.website
  ]
}
