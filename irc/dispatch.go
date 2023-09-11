package irc

import (
	"log"
	"strings"
)

// Dispatch dispatches the events received to registered handlers.
//
// It blocks and return if Disconnect() is called or an error happens, either while
// reading events or from an handler.
func (i *Conn) Dispatch() error {
	for {
		line, err := i.read()
		if err != nil {
			log.Printf("Error received while reading from IRC connexion. Last data received is '%s'.", line)
			return i.errIfConnected(err)
		}

		log.Print(line)
		msg := parseMessage(line)
		var event Event

		switch msg.Command {
		case "PING":
			log.Print("ping pong")
			err := i.sendf("PONG :%s", strings.Join(msg.Parameters, " "))
			if err != nil {
				return i.errIfConnected(err)
			}
		case "001":
			event = RPL_WELCOME
		case "MODE":
			event = MODE
		case "NOTICE":
			event = NOTICE
		case "PRIVMSG":
			event = PRIVMSG
		}

		handler, ok := i.handlers[event]
		if ok {
			err := handler(i, msg)
			if err != nil {
				log.Printf("Error from %v handler: %v", handler, err)
				return err
			}
		}
	}
}

func (i *Conn) errIfConnected(err error) error {
	if i.connected {
		log.Printf("Error while reading from IRC connection: %v", err)
		return err
	}

	return nil
}
