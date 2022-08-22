# Continuous Database Compatbility Testing

This folder contains a few example snippets on how to test database schema downward compatiblity from a CI/CD pipeline. If you want to use this in production you probably want to also use secret management.

## Procedure implemented:
* Fetch current Prod Version
* Copy target database
* Run Integration Tests (Prod Version)
* Apply database migrations for new version
* Deploy app container (Prod Version)
* Run Integration Tests (Prod Version)

## Folders:
* */app* contains a simple Go application that accesses the database
* */cloudbuild* contains an example Cloud Build pipeline that executes the database compatibility check
* */database* contains database migration scripts as SQL statements that get bundled into a docker container image
* */testsuite* contains a chai and mocha based test suite that tests the endpoint of the Golang app

## Setup:
* You will need to create a CloudSQL instane beforehand
* Add the database connection URL (e.g. `postgresql://<username>:<password>@<dbhost>:5432/<dbname>?sslmode=disable`) to the Secret Manager (once for the production database, once for the test database)
* Replace the `PROJECT_ID`, `TEST_DB_URL_SECRET_NAME`, `TEST_DB_URL_SECRET_VERSION`, `PROD_DB_URL_SECRET_NAME`, `PROD_DB_URL_SECRET_VERSION` in the `cloudbuild/app.yaml` file

## Gotchas
* At the moment, the database connection URLs are rather static. This is a problem, when running multiple instances of this database in parallel since the name of the database copy is hardcoded. This is definitely a point of improvement.