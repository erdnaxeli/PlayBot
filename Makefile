export CGO_ENABLED = 0

all: build

build-server:
	go build -ldflags="-s" ./cmd/server

build-ircclient:
	go build -ldflags="-s" ./cmd/ircclient

build: build-ircclient build-server

generate:
	go generate ./...


style:
	go fmt ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1 run

test:
	go test ./...

