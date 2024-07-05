[![deploy](https://github.com/sarumaj/restartable-server/actions/workflows/deploy.yml/badge.svg)](https://github.com/sarumaj/restartable-server/actions/workflows/deploy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sarumaj/restartable-server)](https://goreportcard.com/report/github.com/sarumaj/restartable-server)
[![Maintainability](https://img.shields.io/codeclimate/maintainability-percentage/sarumaj/restartable-server.svg)](https://codeclimate.com/github/sarumaj/restartable-server/maintainability)
[![Go Reference](https://pkg.go.dev/badge/github.com/sarumaj/restartable-server.svg)](https://pkg.go.dev/github.com/sarumaj/restartable-server)
[![Go version](https://img.shields.io/github/go-mod/go-version/sarumaj/restartable-server?logo=go&label=&labelColor=gray)](https://go.dev)

---

# restartable-server

Implementation of a HTTP server resilient to SIGTERM OS signal and restartable on demand.

Especially interesting when deploying HTTP servers with the GO runtime in containers.

Container runtime engine usually send a SIGTERM signal to kill the server process,
which can be catch by the restartable server to restart itself.

Provided as Proof of Concepts. Currently available at [restartable-server.heroku.com](https://restartable-server-a78d4e6a2c84.herokuapp.com/).

https://github.com/sarumaj/restartable-server/assets/71898979/b511aa0e-707c-4260-a8b0-b3179a165e7e

## Local deployment

```
go build -o server cmd/server/main.go
./server
```

## Deployment on Heroku

```
heroku login
docker ps
heroku container:login
heroku container:push web -a <app_name>
heroku container:release web -a <app_name>
```
