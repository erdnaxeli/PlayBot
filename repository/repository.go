package repository

import "github.com/erdnaxeli/PlayBot/types"

type Repository interface {
	// Return a slice of tags for the given music record.
	GetTags(musicRecordId int64) ([]string, error)
	// Save a music post and return the music record id along to a bool which is
	// true if the post is a new one, false is the post already existed. In the
	// latter case, the post is updated.
	SaveMusicPost(types.MusicPost) (int64, bool, error)
	// Save the given tags for the music record pointed by the given id.
	SaveTags(musicRecordId int64, tags []string) error
}
