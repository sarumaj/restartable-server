# restartable-server

Implementation of a HTTP server resilient to SIGTERM OS signal and restartable on demand.

Especially interesting when deploying HTTP servers with the GO runtime in containers.

Container runtime engine usually send a SIGTERM signal to kill the server process,
which can be catch by the restartable server to restart itself.

Provided as Proof of Concepts. Currently available at [restartable-server.apps.ciepiela.org](https://restartable-server.apps.ciepiela.org/).

## Local deployment

```bash
go build -o server cmd/server/main.go
./server
```
