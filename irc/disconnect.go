package irc

// Disconnect sends a QUIT command and close the connection to the IRC server.
func (i *Conn) Disconnect() {
	i.connected = false
	_ = i.sendRaw("QUIT")
	_ = i.conn.Close()
}
