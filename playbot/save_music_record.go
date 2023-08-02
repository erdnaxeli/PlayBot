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
) (int64, types.MusicRecord, error) {
	_, musicRecord, err := p.extractor.Extract(msg)
	if err != nil {
		return 0, types.MusicRecord{}, err
	}

	recordId, err := p.repository.SaveMusicPost(types.MusicPost{
		MusicRecord: musicRecord,
		Person:      person,
		Channel:     channel,
	})
	if err != nil {
		return 0, types.MusicRecord{}, fmt.Errorf("error while saving music record: %w", err)
	}

	return recordId, musicRecord, nil
}
