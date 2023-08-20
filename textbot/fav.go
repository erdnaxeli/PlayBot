package textbot

import (
	"strings"

	"github.com/erdnaxeli/PlayBot/types"
)

// !fav is the only command that can be used with an url. It acts like posting the
// record then executing the !fav command.
//
// Internally, first it parses the message without the "!fav" command, then it executes
// the "!fav" command. The first parameter of the command is optional and is an
// absolute or relative recordID.
func (t *textBot) favCmd(
	channel types.Channel, person types.Person, args []string, user string,
) (Result, error) {
	result, err := t.saveMusicPost(channel, person, strings.Join(args, " "))
	if err != nil {
		return result, err
	}

	if user == "" {
		return result, AuthenticationRequired
	}

	recordID, _, err := t.getRecordIDFromArgs(channel, args)
	if err != nil {
		return result, err
	}

	err = t.playbot.SaveFav(user, recordID)
	return result, err
}
