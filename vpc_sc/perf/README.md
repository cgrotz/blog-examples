


KO_DOCKER_REPO=gcr.io/ipam-autopilot-showcase ko publish --platform all .

export KO_DOCKER_REPO=gcr.io/ipam-autopilot-showcase 

gcloud run deploy backend --quiet \
    --allow-unauthenticated \
    --platform=managed \
    --concurrency=1 \
    --max-instances=10 \
    --region=europe-west1 \
    --image=$(ko publish ./) \
    --args=--processingTime=500 \
    --args=--startupDelay=200

gcloud run deploy layer-1 --quiet \
    --allow-unauthenticated \
    --platform=managed \
    --concurrency=1 \
    --max-instances=10 \
    --region=europe-west1 \
    --image=$(ko publish ./) \
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
    --image=$(ko publish ./) \
    --args=--proxy=true \
    --args=--remote=https://layer-1-pdk3svnohq-ew.a.run.app \
    --args=--preRequestDelay=50 \
    --args=--postRequestDelay=50 \
    --args=--startupDelay=200

