package playbot

import (
	"context"

	"github.com/erdnaxeli/PlayBot/extractors"
	"github.com/erdnaxeli/PlayBot/types"
)

type Repository interface {
	// Return a slice of tags for the given music record.
	GetTags(musicRecordId int64) ([]string, error)
	// Save a music post and return the music record id along to a bool which is
	// true if the post is a new one, false is the post already existed. In the
	// latter case, the post is updated.
	SaveMusicPost(types.MusicPost) (int64, bool, error)
	// Save the given tags for the music record pointed by the given id.
	SaveTags(musicRecordId int64, tags []string) error
	// Search for a music record. It returns a channel to stream SearchResult objects.
	// Closing the channel will produce a panic. If you want to notify than no more
	// results will be needed, cancel the context.
	SearchMusicRecord(
		ctx context.Context, channel types.Channel, words []string, tags []string,
	) (chan SearchResult, error)
}

type SearchResult interface {
	Id() int64
	MusicRecord() types.MusicRecord
}

type Playbot struct {
	extractor  extractors.MultipleSourcesExtractor
	repository Repository
}

func New(extractor extractors.MultipleSourcesExtractor, repository Repository) Playbot {
	return Playbot{
		extractor:  extractor,
		repository: repository,
	}
}
