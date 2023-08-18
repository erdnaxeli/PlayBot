package main

import (
	"log"

	"github.com/erdnaxeli/PlayBot/irc"
)

func main() {
	config := irc.Config{
		Host: "irc.example.org",
		Port: 7000,
		Nick: "playtest",
	}

	conn, err := irc.New(config)
	if err != nil {
		log.Fatal(err)
	}

	conn.OnConnect(func(c *irc.Conn, m irc.Message) error {
		return c.Join("#playbot")
	})
	conn.OnPrivmsg(func(c *irc.Conn, m irc.Message) error {
		log.Printf("Just received a message: %s", m.Parameters[1])
		return nil
	})

	log.Fatal(conn.Dispatch())
}
