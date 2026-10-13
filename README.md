# restartable-server

Implementation of a HTTP server resilient to SIGTERM OS signal and restartable on demand.

Especially interesting when deploying HTTP servers with the GO runtime in containers.

Container runtime engine usually send a SIGTERM signal to kill the server process,
which can be catch by the restartable server to restart itself.

Provided as Proof of Concepts. Currently available at [restartable-server.apps.ciepiela.org](https://restartable-server.apps.ciepiela.org/).

<https://github.com/sarumaj/restartable-server/assets/71898979/b511aa0e-707c-4260-a8b0-b3179a165e7e>

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
