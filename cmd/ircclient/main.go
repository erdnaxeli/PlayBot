package main

import (
	"context"
	"log"
	"net/http"

	"github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/erdnaxeli/PlayBot/config"
	"github.com/erdnaxeli/PlayBot/irc"
)

func main() {
	config, err := config.ReadConfigFile("playbot.conf")
	if err != nil {
		log.Fatalf("Error while reading config file: %s", err)
	}

	ircConfig := irc.Config{
		Host: config.Irc.Host,
		Port: config.Irc.Port,
		Nick: config.Irc.Nick,
	}

	conn, err := irc.New(ircConfig)
	if err != nil {
		log.Fatal(err)
	}

	client := rpc.NewPlaybotCliProtobufClient("http://localhost:1111", &http.Client{})

	conn.OnConnect(func(c *irc.Conn, m irc.Message) error {
		for _, channel := range config.Irc.Channels {
			err := c.Join(channel)
			if err != nil {
				return err
			}
		}

		return nil
	})
	conn.OnPrivmsg(func(c *irc.Conn, m irc.Message) error {
		result, err := client.Execute(
			context.Background(),
			&rpc.TextMessage{
				ChannelName: m.Parameters[0],
				Msg:         m.Parameters[1],
			},
		)
		if err != nil {
			log.Printf("Error while executing command: %s", err)
			return nil
		}

		if result.To == "" || result.Msg == "" {
			return nil
		}

		err = c.Privmsg(result.To, result.Msg)
		if err != nil {
			return err
		}

		return nil
	})

	log.Fatal(conn.Dispatch())
}
