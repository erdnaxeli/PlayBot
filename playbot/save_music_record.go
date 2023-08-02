package playbot

import (
	"fmt"

	"github.com/erdnaxeli/PlayBot/types"
)

// Save a music record pointed by the URL in the given message, and a post of this
// record by the given user in the given channel.
//
// Return the id, the music record, an bool indicating if the record is a new one,
// and an error.
func (p *Playbot) SaveMusicRecord(
	msg string, person types.Person, channel types.Channel,
) (int64, types.MusicRecord, bool, error) {
	_, musicRecord, err := p.extractor.Extract(msg)
	if err != nil {
		return 0, musicRecord, false, err
	}

	recordId, isNew, err := p.repository.SaveMusicPost(types.MusicPost{
		MusicRecord: musicRecord,
		Person:      person,
		Channel:     channel,
	})
	if err != nil {
		return recordId, types.MusicRecord{}, isNew, fmt.Errorf("error while saving music record: %w", err)
	}

	return recordId, musicRecord, isNew, nil
}
