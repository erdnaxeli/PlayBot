package playbot

import (
	"time"

	"github.com/erdnaxeli/PlayBot/types"
)

type MusicRecordStatistics struct {
	PostsCount    int
	PeopleCount   int
	ChannelsCount int

	MaxPerson       types.Person
	MaxPersonCount  int
	MaxChannel      types.Channel
	MaxChannelCount int

	FirstPostPerson  types.Person
	FirstPostChannel types.Channel
	FirstPostDate    time.Time
	FavoritesCount   int
}

func (p *Playbot) GetMusicRecordStatistics(musicRecordID int64) (MusicRecordStatistics, error) {
	return p.repository.GetMusicRecordStatistics(musicRecordID)
}
