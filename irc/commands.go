package irc

func (i *Conn) Privmsg(to string, msg string) error {
	return i.sendf("PRIVMSG %s :%s", to, msg)
}

func (i *Conn) Join(channel string) error {
	return i.sendf("JOIN %s", channel)
}
