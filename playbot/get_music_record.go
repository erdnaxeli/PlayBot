package playbot

import "github.com/erdnaxeli/PlayBot/types"

func (p *Playbot) GetMusicRecord(musicRecordId int64) (types.MusicRecord, error) {
	return p.repository.GetMusicRecord(musicRecordId)
}
