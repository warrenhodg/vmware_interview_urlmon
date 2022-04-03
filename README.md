# urlmon

This application monitors whether certain URLs are
operating correctly, and exposes those results as
Prometheus style metrics.

## Overall application flow

1. Get runtime arguments (`config.go`). This uses `github.com/alecthomas/kong`
   to parse command-line arguments.
2. Setup the HTTP server (`httpserver` package).
3. Setup the healthchecks (`healthcheck` package).
4. Setup the metrics endpoint (`promhttp` package).
5. Setup the checker system (`urlchecker` package). This is a system
  that periodically produces URLs to check, and some workers to
  actually check those URLs.
6. Run the HTTP server.
7. Run the checker system.
8. Wait for signal to shutdown.
9. Fail readiness check for a while.
10. Shutdown the checker system
11. Shutdown the httpserver

## Running the application

The application is able to run either
in standalone mode, but it is recommended that it
rather be run in Kubernetes using the provided helm
chart in the `helm/urlmon` directory.

The application takes several runtime arguments,
which can be seen by passing `--help` when running.

## Developer testing

In order to setup a development environment for
testing the application, a `Makefile` has been provided
which is able to run the application, and various other
applications that would normally already be present in
a production environment.

* To run the usual tests, run `make test`. This runs tests
  while also checking for goroutine race conditions.

* To build a Docker image, run `make docker-build`. It is
  expected that your Docker environment, and command-line
  tooling is already setup.

* To setup a repository into which the container image can
  be pushed (necessary for Kubernetes testing), run
  `make setup-registry`. Also run §`make setup-registry-hostname`
  to ensure that your `docker push` command later can find
  the registry you created.

* To push the built Docker image to your repository, run
  `make docker-push`.

* To setup Prometheus which will scrape the metrics of the
  application, run `make prometheus`.

* To setup Grafana, run `make grafana`. The default username
  and password is `admin`. Setup a datasource to `http://prometheus-server`.
  Future work could precreate this datasource.

* To run the application itself, use `make helm-install`. This will
  run the application in the `urlmon` namespace.

* To be able to browse Grafana from the browser, we use the
  Kubernetes port forwarding to the pod itself, by running
  `make grafana-port-forward`. This allows us to reach
  Grafana at `http://localhost:3000`.

## Notes

* Additional work should be done to allow a version number to be
  specified at build time, and queried at both CLI level, and via
  a `/version` endpoint.

* Additional work should be done to incorporate various log levels
  with the ability to specify the level as a startup argument.

* When this application is incorporated into a formal CI/CD pipeline,
  the container registry should be adjusted.

* Should the number of urls to be checked become too great, then
  a file containing the list would be a better mechanism than command-line
  arguments.
