package main

//go:generate protoc --twirp_out=module=github.com/erdnaxeli/playbot/cmd/cli:. --go_out=module=github.com/erdnaxeli/playbot/cmd/cli:.  service.proto

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/erdnaxeli/PlayBot/cmd/cli/rpc"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "server" {
		log.Fatal(startServer())
	}

	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s CHANNEL PERSON MESSAGE", os.Args[0])
	}

	channel := os.Args[1]
	person := os.Args[2]
	msg := os.Args[3]

	client := rpc.NewPlaybotCliProtobufClient("http://localhost:1111", &http.Client{})
	result, err := client.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channel,
			PersonName:  person,
			Msg:         msg,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Result received")
	fmt.Print(result.Msg)
}
