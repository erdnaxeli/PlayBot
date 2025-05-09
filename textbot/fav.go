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
func (t *TextBot) favCmd(
	channel types.Channel, person types.Person, args []string, user string,
) (Result, error) {
	result, recordID, err := t.saveFavPost(channel, person, args)
	if err != nil {
		return result, err
	}

	if user == "" {
		return result, ErrAuthenticationRequired
	}

	err = t.playbot.SaveFav(user, recordID)
	return result, err
}

func (t *TextBot) saveFavPost(
	channel types.Channel, person types.Person, args []string,
) (Result, int64, error) {
	var result Result
	var recordID int64
	var err error

	if len(args) > 0 {
		result, err = t.saveMusicPost(channel, person, strings.Join(args, " "))
		if err != nil {
			return result, 0, err
		}

		recordID = result.ID
	}

	if recordID == 0 {
		recordID, _, err = t.getRecordIDFromArgs(channel, args)
		if err != nil {
			return Result{}, recordID, err
		}
	}

	return result, recordID, nil
}
