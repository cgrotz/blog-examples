#!/bin/bash
# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

docker build -t gcr.io/ipam-autopilot-showcase/perf:11 .
docker push gcr.io/ipam-autopilot-showcase/perf:11

gcloud run deploy backend --quiet \
    --concurrency=1 \
    --max-instances=1 \
    --min-instances=1 \
    --allow-unauthenticated \
    --platform=managed \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/ipam-autopilot-showcase/perf:11" \
    --args=--processingTime=10

gcloud run deploy layer-1 --quiet \
    --concurrency=1 \
    --max-instances=1 \
    --min-instances=1 \
    --allow-unauthenticated \
    --platform=managed \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/ipam-autopilot-showcase/perf:11" \
    --args=--proxy=true \
    --args=--remote=https://backend-pdk3svnohq-ew.a.run.app

gcloud run deploy layer-1-internal --quiet \
    --concurrency=1 \
    --max-instances=1 \
    --min-instances=1 \
    --allow-unauthenticated \
    --platform=managed \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/ipam-autopilot-showcase/perf:11" \
    --vpc-connector=connector-eu \
    --vpc-egress=all-traffic \
    --args=--proxy=true \
    --args=--remote=https://backend-pdk3svnohq-ew.a.run.app