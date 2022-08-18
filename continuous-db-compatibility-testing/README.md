# Continuous Database Compatbility Testing

This folder contains a few example snippets on how to test database schema downward compatiblity from a CI/CD pipeline. If you want to use this in production you probably want to also use secret management.

Procedure implemented:
* Fetch current Prod Version
* Copy target database
* Run Integration Tests (Prod Version)
* Apply database migrations for new version
* Deploy app container (Prod Version)
* Run Integration Tests (Prod Version)
