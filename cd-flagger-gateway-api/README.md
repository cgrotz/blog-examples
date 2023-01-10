# Canary Deployments with Cloud Deploy, Flagger and Gateway API
```
export GOOGLE_CLOUD_PROJECT_ID=<your_project_on_google_cloud>
export GOOGLE_CLOUD_REGION=<your_google_cloud_region>
```
Create Proxy-only subnet
```
gcloud compute networks subnets create proxy \
    --purpose=REGIONAL_MANAGED_PROXY \
    --role=ACTIVE \
    --region=$GOOGLE_CLOUD_REGION \
    --network=default \
    --range=10.103.0.0/23
```
# GKE Cluster
Cluster with HPA and Workload Identity preinstalled
```
gcloud beta container clusters create "example-cluster" --cluster-version "1.24.5-gke.600" --region "$GOOGLE_CLOUD_REGION"  --machine-type "e2-medium" --max-pods-per-node "30" --num-nodes "1" --enable-autoscaling --min-nodes "0" --max-nodes "3" --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver --enable-managed-prometheus --workload-pool "$GOOGLE_CLOUD_PROJECT_ID.svc.id.goog" --enable-shielded-nodes --gateway-api=standard --enable-ip-alias
```
Connect to cluster
```
gcloud container clusters get-credentials example-cluster --region $GOOGLE_CLOUD_REGION
```

# Bootstrap Flagger and the Gateway
Install Flagger for Gateway API
```
kubectl apply -k github.com/fluxcd/flagger//kustomize/gatewayapi
```

Flagger KSA needs to be annotated, Flagger GSA with monitoring access
```
gcloud iam service-accounts create flagger --project=$GOOGLE_CLOUD_PROJECT_ID

gcloud iam service-accounts add-iam-policy-binding flagger@$GOOGLE_CLOUD_PROJECT_ID.iam.gserviceaccount.com \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:$GOOGLE_CLOUD_PROJECT_ID.svc.id.goog[flagger-system/flagger]"

kubectl annotate serviceaccount flagger \
    --namespace flagger-system \
    iam.gke.io/gcp-service-account=flagger@$GOOGLE_CLOUD_PROJECT_ID.iam.gserviceaccount.com
```
Create certificate for Gateway
```
gcloud compute ssl-certificates create app-dev-grotz-dev \
    --domains=app.dev.grotz.dev \
    --global

gcloud compute ssl-certificates create app-prod-grotz-prod \
    --domains=app.prod.grotz.dev \
    --global
```

Deploy App
```
skaffold run --default-repo=gcr.io/$GOOGLE_CLOUD_PROJECT_ID
```
Fetch IP for DNS setup
```
kubectl get gateways.gateway.networking.k8s.io app-dev  -n dev -o=jsonpath="{.status.addresses[0].value}"
```
```
kubectl apply -f bootstrap.yaml
```
Deploy GMP Query Interface
```
kubectl create serviceaccount gmp -n prod
gcloud iam service-accounts create gmp-sa --project=$GOOGLE_CLOUD_PROJECT_ID

gcloud iam service-accounts add-iam-policy-binding gmp-sa@$GOOGLE_CLOUD_PROJECT_ID.iam.gserviceaccount.com \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:$GOOGLE_CLOUD_PROJECT_ID.svc.id.goog[prod/gmp]"

gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT_ID \
  --member=serviceAccount:gmp-sa@$GOOGLE_CLOUD_PROJECT_ID.iam.gserviceaccount.com \
  --role=roles/monitoring.viewer

kubectl annotate serviceaccount gmp \
    --namespace prod \
    iam.gke.io/gcp-service-account=gmp-sa@$GOOGLE_CLOUD_PROJECT_ID.iam.gserviceaccount.com

kubectl apply -n prod -f gmp-frontend.yaml
```

# Create Cloud Deploy Pipeline
Set permissions for Cloud Deploy and apply pipeline

gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT_ID \
    --member=serviceAccount:$(gcloud projects describe $GOOGLE_CLOUD_PROJECT_ID \
    --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
    --role="roles/clouddeploy.jobRunner"
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT_ID \
    --member=serviceAccount:$(gcloud projects describe $GOOGLE_CLOUD_PROJECT_ID \
    --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
    --role="roles/container.developer"

gcloud deploy apply --file clouddeploy.yaml --region=$GOOGLE_CLOUD_REGION --project=$GOOGLE_CLOUD_PROJECT_ID 

# Trigger Pipeline
Create new release for deployment
 skaffold run --default-repo=gcr.io/$GOOGLE_CLOUD_PROJECT_ID -p prod

skaffold build --default-repo=gcr.io/$GOOGLE_CLOUD_PROJECT_ID 

gcloud deploy releases create release-001 \
  --project=$GOOGLE_CLOUD_PROJECT_ID  \
  --region=$GOOGLE_CLOUD_REGION \
  --delivery-pipeline=canary \
  --images=skaffold-kustomize=gcr.io/$GOOGLE_CLOUD_PROJECT_ID/skaffold-kustomize:899d24a-dirty

curl --header 'Host: app.dev.grotz.dev' http://10.132.0.48
curl --header 'Host: app.prod.grotz.dev' http://10.132.0.48

After the release is deploy to dev, promote it to production
gcloud deploy releases promote  --release=release-001 --delivery-pipeline=canary --region=$GOOGLE_CLOUD_REGION --to-target=prod

## Observe the pipeline
kubectl -n prod describe canary/app
gcloud compute url-maps export gkegw1-ll1w-prod-app-4bpekl57o1qy --region=$GOOGLE_CLOUD_REGION
