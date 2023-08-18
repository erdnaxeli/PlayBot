package irc

func (i *Conn) Disconnect() {
	i.connected = false
	_ = i.sendRaw("QUIT")
	_ = i.conn.Close()
}
