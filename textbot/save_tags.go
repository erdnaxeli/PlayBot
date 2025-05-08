package textbot

import (
	"github.com/erdnaxeli/PlayBot/types"
)

func (t *textBot) saveTagsCmd(
	channel types.Channel, _ types.Person, args []string,
) error {
	if len(args) == 0 {
		return nil
	}

	recordID, args, err := t.getRecordIDFromArgs(channel, args)
	if err != nil {
		return err
	}

	removeHash(args)
	err = t.playbot.SaveTags(recordID, args)
	return err
}

func removeHash(words []string) {
	for idx, word := range words {
		if word[0] == '#' {
			words[idx] = word[1:]
		}
	}
}
