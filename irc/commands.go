package irc

func (i *Conn) Join(channel string) error {
	return i.sendf("JOIN %s", channel)
}

func (i *Conn) Mode(target string, modes string) error {
	return i.sendf("MODE %s %s", target, modes)
}

func (i *Conn) Privmsg(to string, msg string) error {
	return i.sendf("PRIVMSG %s :%s", to, msg)
}
