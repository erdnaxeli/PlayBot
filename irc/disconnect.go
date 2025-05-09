package irc

// Disconnect sends a QUIT command and close the connection to the IRC server.
//
// I
func (i *Conn) Disconnect() {
	i.connected = false
	_ = i.sendf("QUIT")
	_ = i.conn.Close()
}
