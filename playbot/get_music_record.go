package playbot

import "github.com/erdnaxeli/PlayBot/types"

// GetMusicRecord returns a MusicRecord object.
func (p *Playbot) GetMusicRecord(musicRecordID int64) (types.MusicRecord, error) {
	return p.repository.GetMusicRecord(musicRecordID)
}
