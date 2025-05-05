package playbot

import "github.com/erdnaxeli/PlayBot/types"

// SaveMusicPost create a post for the given music record, by the given person, on the given channel.
//
// This is useful for example if your app allows a user to share an existing record to somewhere else.
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
