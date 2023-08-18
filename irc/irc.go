package irc

import (
	"fmt"
	"net"
	"net/textproto"
	"time"
)

type Config struct {
	Host string
	Port int
	Nick string
}

type Conn struct {
	config           Config
	connected        bool
	handlers         map[Event]Handler
	throttleDeadline time.Time

	conn   net.Conn
	reader *textproto.Reader
	writer *textproto.Writer
}

// Connect with the given config to the server and return a new IrcConnection. It uses
// tls for connection.
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
	err := i.writer.PrintfLine(msg)
	i.throttleDeadline = i.throttleDeadline.Add(500 * time.Millisecond)
	return err
}

func (i *Conn) sendf(format string, args ...any) error {
	return i.sendRaw(fmt.Sprintf(format, args...))
}

func (i *Conn) read() (string, error) {
	return i.reader.ReadLine()
}
