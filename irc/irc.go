// Package irc implements the IRC protocol.
//
// It allows to connect to a server, react to received events, and send commands.
package irc

import (
	"fmt"
	"log"
	"net"
	"net/textproto"
	"time"
)

// Config is the configuration to connect to an IRC server.
type Config struct {
	Host string
	Port int
	Nick string
}

// Conn represents an IRC connection.
//
// The zero value is not usable, you should use New().
type Conn struct {
	config           Config
	connected        bool
	handlers         map[Event]Handler
	throttleDeadline time.Time

	conn   net.Conn
	reader *textproto.Reader
	writer *textproto.Writer
}

// New connects with the given config to the server and return a new IrcConnection.
//
// It uses tls for connection.
// The connection is buffered for reads and writes.
func New(config Config) (*Conn, error) {
	irc := &Conn{
		config:    config,
		connected: false,
		handlers:  make(map[Event]Handler),
	}
	err := irc.Connect()
	if err != nil {
		return nil, err
	}

	return irc, nil
}

func (i *Conn) sendRaw(msg string) error {
	time.Sleep(time.Until(i.throttleDeadline))
	err := i.writer.PrintfLine("%s", msg)
	i.throttleDeadline = time.Now().Add(500 * time.Millisecond)

	if err != nil {
		log.Printf("Error while sending \"%s\": %v", msg, err)
	}
	return err
}

func (i *Conn) sendf(format string, args ...any) error {
	return i.sendRaw(fmt.Sprintf(format, args...))
}

func (i *Conn) read() (string, error) {
	return i.reader.ReadLine()
}
