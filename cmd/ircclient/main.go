package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/erdnaxeli/PlayBot/config"
	"github.com/erdnaxeli/PlayBot/irc"
)

type bot struct {
	client rpc.PlaybotCli
	config config.Config
}

func main() {
	fmt.Print("Starting IRC client...")
	config, err := config.ReadConfig()
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

	client := rpc.NewPlaybotCliProtobufClient(
		fmt.Sprintf("http://%s", config.ServerAddress), &http.Client{},
	)
	b := bot{client, config}

	conn.OnConnect(b.onConnect)
	conn.OnNotice(b.onMessage)
	conn.OnPrivMsg(b.onMessage)

	log.Fatal(conn.Dispatch())
}

func (b bot) onConnect(c *irc.Conn, m irc.Message) error {
	err := c.Privmsg("NickServ", "identify "+b.config.Irc.NickServPassword)
	if err != nil {
		return err
	}

	for _, channel := range b.config.Irc.Channels {
		err := c.Join(channel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b bot) onMessage(c *irc.Conn, m irc.Message) error {
	nick := c.GetNick(m.Prefix)
	if nick == "" {
		return nil
	}

	b.exec(
		&rpc.TextMessage{
			ChannelName: m.Parameters[0],
			Msg:         m.Parameters[1],
			PersonName:  nick,
		},
		c,
	)

	return nil
}

func (b bot) exec(msg *rpc.TextMessage, c *irc.Conn) {
	result, err := b.client.Execute(context.Background(), msg)
	if err != nil {
		log.Printf("Error while executing command: %s", err)
		return
	}

	for _, msg := range result.Msg {
		log.Print(msg)
		err = c.Privmsg(msg.To, msg.Msg)
		if err != nil {
			return
		}
	}
}
