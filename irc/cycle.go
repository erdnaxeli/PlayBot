package irc

func (i Conn) Cycle() error {
	i.Disconnect()
	return i.Connect()
}
