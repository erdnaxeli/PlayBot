package main

//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
//go:generate go install github.com/twitchtv/twirp/protoc-gen-twirp@v8.1.3
//go:generate protoc --twirp_out=module=github.com/erdnaxeli/playbot/cmd/cli:. --go_out=module=github.com/erdnaxeli/playbot/cmd/cli:.  service.proto

import (
	"log"
)

func main() {
	log.Fatal(startServer())
}
