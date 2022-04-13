


KO_DOCKER_REPO=gcr.io/ipam-autopilot-showcase ko publish --platform all .

export KO_DOCKER_REPO=gcr.io/ipam-autopilot-showcase 

docker build -t gcr.io/ipam-autopilot-showcase/perf:2 .
docker push gcr.io/ipam-autopilot-showcase/perf:2

gcloud run deploy backend --quiet \
    --allow-unauthenticated \
    --platform=managed \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/ipam-autopilot-showcase/perf:2" \
    --args=--processingTime=100

gcloud run deploy layer-1 --quiet \
    --allow-unauthenticated \
    --platform=managed \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/ipam-autopilot-showcase/perf:2" \
    --args=--proxy=true \
    --args=--remote=https://backend-pdk3svnohq-ew.a.run.app


gcloud run deploy layer-1-internal --quiet \
    --allow-unauthenticated \
    --platform=managed \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/ipam-autopilot-showcase/perf:2" \
    --vpc-connector=connector \
    --vpc-egress=all-traffic \
    --args=--proxy=true \
    --args=--remote=https://backend-pdk3svnohq-ew.a.run.app



gcloud run deploy backend --quiet \
    --allow-unauthenticated \
    --platform=managed \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/ipam-autopilot-showcase/perf:2" \
    --args=--processingTime=100 \
    --args=--startupDelay=200

gcloud run deploy layer-1 --quiet \
    --allow-unauthenticated \
    --platform=managed \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/ipam-autopilot-showcase/perf:2" \
    --args=--proxy=true \
    --args=--remote=https://backend-pdk3svnohq-ew.a.run.app \
    --args=--preRequestDelay=50 \
    --args=--postRequestDelay=50 \
    --args=--startupDelay=200

gcloud run deploy layer-2 --quiet \
    --allow-unauthenticated \
    --platform=managed \
    --concurrency=1 \
    --max-instances=10 \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/ipam-autopilot-showcase/perf:1" \
    --args=--proxy=true \
    --args=--remote=https://layer-1-pdk3svnohq-ew.a.run.app \
    --args=--preRequestDelay=50 \
    --args=--postRequestDelay=50 \
    --args=--startupDelay=200

