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
		case "PRIVMSG":
			event = PRIVMSG
		}

		handler, ok := i.handlers[event]
		if ok {
			err := handler(i, msg)
			if err != nil {
				return err
			}
		}
	}
}

func (i *Conn) errIfConnected(err error) error {
	if i.connected {
		return err
	}

	return nil
}
