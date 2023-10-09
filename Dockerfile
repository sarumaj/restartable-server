FROM golang:latest

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build \
    -trimpath \
    -ldflags="-s -w -extldflags=-static" \
    -tags="osusergo netgo static_build" \
    -o "/restartable-server" \
    "cmd/server/main.go" && rm -rf .

CMD ["/restartable-server"]