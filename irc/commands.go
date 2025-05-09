package irc

// Join sends a JOIN command.
func (i *Conn) Join(channel string) error {
	return i.sendf("JOIN %s", channel)
}

// Mode sends a MODE command.
func (i *Conn) Mode(target string, modes string) error {
	return i.sendf("MODE %s %s", target, modes)
}

// Privmsg sends a PRIVMSG command.
func (i *Conn) Privmsg(to string, msg string) error {
	return i.sendf("PRIVMSG %s :%s", to, msg)
}
