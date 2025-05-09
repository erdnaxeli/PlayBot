// Package irc implements the IRC protocol.
//
// It allows to connect to a server, react to received events, and send commands.
package irc

import (
	"bufio"
	"crypto/tls"
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
//
// The received events are not read until you call the Dispatch(), not even the PING event, so you should call it quickly to avoid any server timeout.
func New(config Config) (*Conn, error) {
	conn, err := tls.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to dial tls connection: %w", err)
	}

	return NewWithConn(config, conn)
}

// NewWithConn create a Conn object using the given connection.
//
// The received events are not read until you call the Dispatch(), not even the PING event, so you should call it quickly to avoid any server timeout.
func NewWithConn(config Config, conn net.Conn) (*Conn, error) {
	irc := &Conn{
		config:    config,
		connected: false,
		handlers:  make(map[Event]Handler),

		conn:   conn,
		reader: textproto.NewReader(bufio.NewReader(conn)),
		writer: textproto.NewWriter(bufio.NewWriter(conn)),
	}

	err := irc.Connect()
	if err != nil {
		return nil, err
	}

	return irc, nil
}

func (i *Conn) sendf(format string, args ...any) error {
	time.Sleep(time.Until(i.throttleDeadline))
	err := i.writer.PrintfLine(format, args...)
	i.throttleDeadline = time.Now().Add(500 * time.Millisecond)

	if err != nil {
		log.Printf("Error while sending \"%s\": %v", format, err)
	}
	return err
}

func (i *Conn) read() (string, error) {
	return i.reader.ReadLine()
}
