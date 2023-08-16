package textbot

import (
	"log"
	"regexp"

	"github.com/erdnaxeli/PlayBot/types"
)

func (t *textBot) saveMusicPost(
	channel types.Channel, person types.Person, msg string,
) (Result, error) {
	recordID, musicRecord, isNew, err := t.playbot.ParseAndSaveMusicRecord(
		msg, person, channel,
	)
	if err != nil {
		return Result{}, err
	}

	log.Println("record saved", recordID, musicRecord)
	err = t.saveTags(msg, recordID)
	if err != nil {
		return Result{}, err
	}

	tags, err := t.playbot.GetTags(recordID)
	if err != nil {
		return Result{}, err
	}

	result := Result{
		ID:          recordID,
		MusicRecord: musicRecord,
		Tags:        tags,
		IsNew:       isNew,
	}
	return result, nil
}

func (t *textBot) saveTags(msg string, recordId int64) error {
	tags := extractTagsFromMsg(msg)
	err := t.playbot.SaveTags(recordId, tags)
	return err
}

func extractTagsFromMsg(msg string) []string {
	re := regexp.MustCompile(`\s+`)
	args := re.Split(msg, -1)
	return extractTags(args)
}

func extractTags(args []string) []string {
	var tags []string
	for _, word := range args {
		if len(word) > 0 && word[0] == '#' {
			tags = append(tags, word[1:])
		}
	}

	return tags
}
