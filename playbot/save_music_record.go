package playbot

import "github.com/erdnaxeli/PlayBot/types"

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

	recordId, err := p.repository.SaveMusicRecord(types.MusicPost{
		MusicRecord: musicRecord,
		Person:      person,
		Channel:     channel,
	})
	if err != nil {
		return 0, err
	}

	return recordId, nil
}

func (p *Playbot) SaveTags(musicRecordId int64, tags []string) error {
	return p.repository.SaveTags(musicRecordId, tags)
}
