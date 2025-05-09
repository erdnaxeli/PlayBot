package irc

// Connect connect to the IRC server.
//
// It starts the TCP connection, and then send a NICK command followed by an
// USER command.
func (i *Conn) Connect() error {
	err := i.sendf("NICK %s", i.config.Nick)
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
