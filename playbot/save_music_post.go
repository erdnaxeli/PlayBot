package playbot

import "github.com/erdnaxeli/PlayBot/types"

func (p *Playbot) SaveMusicPost(
	recordID int64, channel types.Channel, person types.Person,
) error {
	record, err := p.repository.GetMusicRecord(recordID)
	if err != nil {
		return err
	}

	_, _, err = p.repository.SaveMusicPost(types.MusicPost{
		MusicRecord: record,
		Person:      person,
		Channel:     channel,
	})
	if err != nil {
		return err
	}

	return nil
}
