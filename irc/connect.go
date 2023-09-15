package irc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/textproto"
)

func (i *Conn) Connect() error {
	var err error
	i.conn, err = tls.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", i.config.Host, i.config.Port),
		nil,
	)
	if err != nil {
		log.Printf("Unable to dial tls connection: %v", err)
		return err
	}

	i.reader = textproto.NewReader(bufio.NewReader(i.conn))
	i.writer = textproto.NewWriter(bufio.NewWriter(i.conn))

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
