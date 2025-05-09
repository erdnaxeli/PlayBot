// Package irc implements the IRC protocol.
//
// It allows to connect to a server, react to received events, and send commands.
package irc

import (
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

	// SocketFactory is used to open the connection to the IRC server.
	//
	// The factory may be called multiple times if Cycle() is used.
	SocketFactory func(Config) (net.Conn, error)
}

// TLSConfig returns a configuration object using a TLS connection.
func TLSConfig(host string, port int, nick string) Config {
	return Config{
		Host: host,
		Port: port,
		Nick: nick,

		SocketFactory: func(config Config) (net.Conn, error) {
			conn, err := tls.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", config.Host, config.Port),
				nil,
			)
			if err != nil {
				return nil, fmt.Errorf("unable to dial tls connection: %w", err)
			}

			return conn, nil
		},
	}
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

// New connects with the given config to the server and return a new Conn.
//
// It sends connection IRC commands right away.
//
// The connection is buffered for reads and writes.
//
// The received events are not read until you call the Dispatch(), not even the PING event, so you should call it quickly to avoid any server timeout.
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
