export CGO_ENABLED = 0

GO = go1.21.0

all: build

build-server:
	$(GO) build -ldflags="-s -w" ./cmd/server

build-ircclient:
	$(GO) build -ldflags="-s -w" ./cmd/ircclient

build: build-ircclient build-server

generate:
	$(GO) generate ./...


style:
	$(GO) fmt ./...
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1 run

test:
	$(GO) test ./...

