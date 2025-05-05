package textbot

import "github.com/erdnaxeli/PlayBot/types"

func (t *textBot) statsCmd(
	channel types.Channel, _ types.Person, args []string,
) (Result, error) {
	if len(args) > 1 {
		return Result{}, ErrInvalidUsage
	}

	recordID, _, err := t.getRecordIDFromArgs(channel, args)
	if err != nil {
		return Result{}, err
	}

	stats, err := t.playbot.GetMusicRecordStatistics(recordID)
	if err != nil {
		return Result{}, err
	}

	return Result{Statistics: stats}, nil
}
