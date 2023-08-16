package textbot

import (
	"errors"
	"strconv"

	"github.com/erdnaxeli/PlayBot/types"
)

var OffsetToBigError = errors.New("offset too big")

func (t *textBot) saveTagsCmd(
	channel types.Channel, person types.Person, args []string,
) error {
	if len(args) == 0 {
		return nil
	}

	recordID, args := parseID(args)
	if recordID <= 0 {
		if recordID < -10 {
			return OffsetToBigError
		}

		var err error
		recordID, err = t.playbot.GetLastID(channel, int(recordID))
		if err != nil {
			return err
		}
	}

	tags := extractTags(args)
	err := t.playbot.SaveTags(recordID, tags)
	return err
}

func parseID(args []string) (int64, []string) {
	recordID, err := strconv.ParseInt(args[0], 10, 64)
	if err == nil {
		return recordID, args[1:]
	}

	return 0, args
}
