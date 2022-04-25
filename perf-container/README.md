# Perf Container

This is a very simple container for Google Cloud Run that I use to measure different performance metrics.

You can run it in two modes:
* *backend* which responds to all requests by printing the request object to the response
* *proxy* which works as a reverse proxy to a configurable remote

You can configure various behaviors of the container to analyze the performance of your platform by provding them as args to the container:
| Variable         | Env                       | Mode    | Default | Description                                                                        | Example Value              |
|------------------|---------------------------|---------|---------|------------------------------------------------------------------------------------|----------------------------|
| preRequestDelay  | PRE_REQUEST_DELAY         | proxy   | 0       | Delay in Milliseconds before the request is passed to the backend                  | dev, qa, prod              |
| postRequestDelay | POST_REQUEST_DELAY        | proxy   | 0       | Delay in Milliseconds after the request is passed to the backend                   | sa@<project>.landisgyr.com |
| processingTime   | PROCESSING_TIME           | backend | 0       | Time in Milliseconds before for fake processing before the request is responded to |                            |
| startupDelay     | STARTUP_DELAY             | both    | 0       | Delay in Milliseconds before the container starts up                               |                            |
| remote           | REVERSE_PROXY_DESTINATION | proxy   |         | Backend for the reverse proxy                                                      |                            |
| proxy            | RUN_AS_REVERSE_PROXY      | proxy   | false   | Should the app run in reverse proxy mode                                           |                            |
| port             | PORT                      | N/A     | 8080    | Server port for the app                                                            |                            |
| tracing          | TRACING                   | N/A     | true    | Tracing enabled; defaults to true                                                  |                            |


Example deployment to Cloud Run in backend mode:
```
gcloud run deploy backend --quiet \
    --allow-unauthenticated \
    --platform=managed \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/<project_id>/perf:1" \
    --args=--processingTime=100
```

Example deployment to Cloud Run in proxy mode:
```
gcloud run deploy frontend --quiet \
    --allow-unauthenticated \
    --platform=managed \
    --region=europe-west1 \
    --project=ipam-autopilot-showcase \
    --image="gcr.io/<project_id>/perf:1" \
    --args=--proxy=true \
    --args=--remote=https://<proxied_service>.a.run.app
```

The container writes traces to Cloud Tracing and produces 3 signals using structured logs:
A started event:
```
{
    ...
    "jsonPayload": {
        "type":  "event",
        "event": "STARTED"
    }
    ...
}
```

A stopped event:
```
{
    ...
    "jsonPayload": {
        "type":  "event",
        "event": "STOPPED"
    }
    ...
}
```
A `proxy_time` event that measures the time it takes the proxy container to receive a response from the server, that you can extract as a [log based metric](https://cloud.google.com/logging/docs/logs-based-metrics):
```
{
    ...
    "jsonPayload": {
        "proxy_time": 100
    }
    ...
}
```

