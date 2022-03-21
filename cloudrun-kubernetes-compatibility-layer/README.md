# Cloud Run Kubernetes Compatibility Layer

This is a simple application that is intended to be run on Google Cloud Run. It implements a reverse proxy for the Cloud Run API. You can specifiy which regional endpoint you want to front with the API using the CLOUD_RUN_API_HOST envrionment variable. The app also implements several Kubernetes discovery endpoints (e.g. `/api`, `/api/v1`, `/version`, `/apis` and `/apis/serving.knative.dev/v1`)

The Cloud Run service currently requires that you allow `run.invoker` to allUsers. It has a passthrough for the Auth Token, provided in the `Authorization`-Header.

If you have [ko](https://github.com/google/ko) and the [gcloud sdk](https://cloud.google.com/sdk/docs/install) installed, you can simple deploy the application from this folder using.
```
gcloud run deploy cr-k8s-compatibility --image $(KO_DOCKER_REPO=gcr.io/<my-project> ko build ./) --allow-unauthenticated
```
The result should print the URL of the Cloud Run service, use this as cluster endpoint in your `.kube/config`. Following is an example Kubernetes config:
```
apiVersion: v1
clusters:
- cluster:
    server: https://<your cloud run service>:443
  name: cloudrun
contexts:
- context:
    cluster: cloudrun
    user: cloudrun
  name: cloudrun
current-context: cloudrun
kind: Config
preferences: {}
users:
- name: cloudrun
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      args: null
      command: gke-gcloud-auth-plugin
      env: null
      installHint: Install gke-gcloud-auth-plugin for use with kubectl by following
        go/gke-kubectl-exec-auth
      interactiveMode: IfAvailable
      provideClusterInfo: true
```