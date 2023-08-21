package playbot

import (
	"context"
	"sync"

	"github.com/erdnaxeli/PlayBot/types"
)

type Repository interface {
	// Return the ID of the last music record in the given channel.
	// If an offset is given, that many music record are skipped before selecting the
	// to return the id from.
	GetLastID(channel types.Channel, offset int) (int64, error)
	// Get a music record by ID.
	GetMusicRecord(musicRecordID int64) (types.MusicRecord, error)
	GetMusicRecordStatistics(musicRecordID int64) (MusicRecordStatistics, error)
	// Return a slice of tags for the given music record.
	GetTags(musicRecordId int64) ([]string, error)
	// Save the given record into the user's favorites.
	//
	// If the record was already in the user's favorites, return no error. If the record
	// does not exist, return a NoRecordFoundError.
	SaveFav(user string, recordID int64) error
	// Save a music post and return the music record id along to a bool which is
	// true if the post is a new one, false is the post already existed. In the
	// latter case, the post is updated.
	SaveMusicPost(types.MusicPost) (int64, bool, error)
	// Save the given tags for the music record pointed by the given ID.
	SaveTags(musicRecordId int64, tags []string) error
	// Search for a music record. It returns a channel to stream SearchResult objects.
	SearchMusicRecord(
		ctx context.Context, channel types.Channel, words []string, tags []string,
	) (int64, chan SearchResult, error)
}

type Extractor interface {
	// Given an URL, Extract tries to extract the record information and returns it.
	Extract(url string) (types.MusicRecord, error)
}

type SearchResult interface {
	Id() int64
	MusicRecord() types.MusicRecord
}

type Playbot struct {
	extractor  Extractor
	repository Repository

	// Contains the ongoing searches.
	searches      map[types.Channel]searchCursor
	searchesMutex sync.RWMutex
}

type searchCursor struct {
	cancel func()
	count  int64
	ch     chan SearchResult
	search Search
}

func New(extractor Extractor, repository Repository) *Playbot {
	return &Playbot{
		extractor:  extractor,
		repository: repository,
		searches:   make(map[types.Channel]searchCursor),
	}
}
