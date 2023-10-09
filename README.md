[![Go Report Card](https://goreportcard.com/badge/github.com/sarumaj/restartable-server)](https://goreportcard.com/report/github.com/sarumaj/restartable-server)
[![Maintainability](https://img.shields.io/codeclimate/maintainability-percentage/sarumaj/restartable-server.svg)](https://codeclimate.com/github/sarumaj/restartable-server/maintainability)

---

# restartable-server

Implementation of a HTTP server resilient to SIGTERM OS signal and restartable on demand.

Provided as Proof of Concepts. Currently available at [restartable-server.heroku.com](https://restartable-server-a78d4e6a2c84.herokuapp.com/).

## Local deployment

```
go -o server build cmd/server/main.go
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
