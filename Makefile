export CGO_ENABLED = 0

all: build

build-server:
	go build -ldflags="-s" ./cmd/server

build-ircclient:
	go build -ldflags="-s" ./cmd/ircclient

build: build-ircclient build-server

generate:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install github.com/twitchtv/twirp/protoc-gen-twirp
	protoc --twirp_out=module=github.com/erdnaxeli/playbot/cmd/cli:cmd/server/ --go_out=module=github.com/erdnaxeli/playbot/cmd/cli:cmd/server/ cmd/server/service.proto

style:
	go fmt ./...
	golangci-lint run --fix ./...

test:
	go test ./...

