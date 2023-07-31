package playbot

import (
	"fmt"

	"github.com/erdnaxeli/PlayBot/types"
)

// Save a music record pointed by the URL in the given message.
//
// Return the matched URL (if any), and an error.
func (p *Playbot) SaveMusicRecord(
	msg string, person types.Person, channel types.Channel,
) (int64, error) {
	_, musicRecord, err := p.extractor.Extract(msg)
	if err != nil {
		return 0, err
	}

	recordId, err := p.repository.SaveMusicPost(types.MusicPost{
		MusicRecord: musicRecord,
		Person:      person,
		Channel:     channel,
	})
	if err != nil {
		return 0, fmt.Errorf("error while saving music record: %w", err)
	}

	return recordId, nil
}

func (p *Playbot) SaveTags(musicRecordId int64, tags []string) error {
	err := p.repository.SaveTags(musicRecordId, tags)
	if err != nil {
		return fmt.Errorf("error while saving tags: %w", err)
	}

	return nil
}
