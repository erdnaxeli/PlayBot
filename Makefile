export CGO_ENABLED = 0

generate:
	go generate ./...

build-cli:
	go build -ldflags="-s -w" ./cmd/cli

style:
	go fmt ./...
	golangci-lint run

test:
	go test ./...

