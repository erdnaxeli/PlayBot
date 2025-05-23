package irc

import "strings"

// Message represents a parsed IRC message.
type Message struct {
	Prefix     string
	Command    string
	Parameters []string
}

func parseMessage(msg string) Message {
	var message Message
	parts := strings.Split(msg, " ")
	if len(parts) == 0 {
		return message
	}

	if parts[0][0] == ':' {
		message.Prefix = parts[0][1:]
		parts = parts[1:]
	}

	message.Command = parts[0]
	message.Parameters = parseParameters(parts[1:])

	return message
}

func parseParameters(parts []string) []string {
	var parameters []string
	for i, part := range parts {
		if len(part) == 0 {
			continue
		}

		if part[0] != ':' {
			parameters = append(parameters, part)
			continue
		}

		// If the parameter start with ":", it means all that follow is part of the same parameter.
		parts[i] = part[1:]
		parameters = append(parameters, strings.Join(parts[i:], " "))
		break
	}

	return parameters
}

// GetNick extracts the nick part from a prefix in the form "nick!username@host".
// If the prefix is invalid, it returns an empty string.
func (i *Conn) GetNick(prefix string) string {
	nick, _, found := strings.Cut(prefix, "!")
	if !found {
		return ""
	}

	return nick
}
