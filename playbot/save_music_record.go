package playbot

import (
	"fmt"

	"github.com/erdnaxeli/PlayBot/types"
)

// ParseAndSaveMusicRecord save a music record pointed by the URL in the given message,
// and create post for this record by the given user in the given channel.
//
// Return the id, the music record, an bool indicating if the record is a new one,
// and an error.
func (p *Playbot) ParseAndSaveMusicRecord(
	msg string, person types.Person, channel types.Channel,
) (recordID int64, record types.MusicRecord, isNew bool, err error) {
	musicRecord, err := p.extractor.Extract(msg)
	if err != nil {
		return 0, types.MusicRecord{}, false, err
	}

	if (musicRecord == types.MusicRecord{}) {
		return 0, musicRecord, false, nil
	}

	recordID, isNew, err = p.repository.SaveMusicPost(types.MusicPost{
		MusicRecord: musicRecord,
		Person:      person,
		Channel:     channel,
	})
	if err != nil {
		return recordID, types.MusicRecord{}, isNew, fmt.Errorf("error while saving music record: %w", err)
	}

	return recordID, musicRecord, isNew, nil
}
