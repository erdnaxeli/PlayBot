package irc

// Disconnect sends a QUIT command and close the connection to the IRC server.
//
// This method ignores any error.
// It remove the connection to the IRC server, meaning any command sent after that
// while fails except Connect().
func (i *Conn) Disconnect() {
	i.connected = false
	_ = i.sendf("QUIT")
	_ = i.conn.Close()
	i.conn = nil
}
