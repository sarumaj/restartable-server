FROM golang:latest AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go test -v ./... && \
    go build \
    -trimpath \
    -ldflags="-s -w -extldflags=-static" \
    -tags="osusergo netgo static_build" \
    -o /server \
    "cmd/server/main.go" && \
    rm -rf /usr/src/app

CMD ["/server"]