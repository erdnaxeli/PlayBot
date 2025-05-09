package irc

import (
	"bufio"
	"fmt"
	"net/textproto"
)

// Connect connect to the IRC server.
//
// It starts the TCP connection, and then send a NICK command followed by an
// USER command.
func (i *Conn) Connect() error {
	conn, err := i.config.SocketFactory(i.config)
	if err != nil {
		return fmt.Errorf("error while creating the new connection using the provided SocketFactory: %w", err)
	}
	i.conn = conn

	i.reader = textproto.NewReader(bufio.NewReader(conn))
	i.writer = textproto.NewWriter(bufio.NewWriter(conn))

	err = i.sendf("NICK %s", i.config.Nick)
	if err != nil {
		return err
	}

	err = i.sendf("USER %s 0 * :%s", i.config.Nick, i.config.Nick)
	if err != nil {
		return err
	}

	i.connected = true
	return nil
}
