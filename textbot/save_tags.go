package textbot

import (
	"errors"

	"github.com/erdnaxeli/PlayBot/types"
)

var OffsetToBigError = errors.New("offset too big")

func (t *textBot) saveTagsCmd(
	channel types.Channel, person types.Person, args []string,
) error {
	if len(args) == 0 {
		return nil
	}

	recordID, args, err := t.getRecordIDFromArgs(channel, args)
	if err != nil {
		return err
	}

	tags := extractTags(args)
	err = t.playbot.SaveTags(recordID, tags)
	return err
}
