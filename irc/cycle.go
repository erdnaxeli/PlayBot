package irc

// Cycle disconnects from the IRC server and then connect again.
func (i Conn) Cycle() error {
	i.Disconnect()
	return i.Connect()
}
