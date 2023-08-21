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

	recordId, isNew, err := p.repository.SaveMusicPost(types.MusicPost{
		MusicRecord: musicRecord,
		Person:      person,
		Channel:     channel,
	})
	if err != nil {
		return recordId, types.MusicRecord{}, isNew, fmt.Errorf("error while saving music record: %v", err)
	}

	return recordId, musicRecord, isNew, nil
}
