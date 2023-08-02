package repository

import "github.com/erdnaxeli/PlayBot/types"

type Repository interface {
	// Return a slice of tags for the given music record.
	GetTags(musicRecordId int64) ([]string, error)
	// Save a music post and return the music record id.
	SaveMusicPost(types.MusicPost) (int64, error)
	// Save the given tags for the music record pointed by the given id.
	SaveTags(musicRecordId int64, tags []string) error
}
