package playbot

import (
	"github.com/erdnaxeli/PlayBot/extractors"
	"github.com/erdnaxeli/PlayBot/repository"
)

type Playbot struct {
	extractor  extractors.MultipleSourcesExtractor
	repository repository.Repository
}

func New(extractor extractors.MultipleSourcesExtractor, repository repository.Repository) *Playbot {
	return &Playbot{
		extractor:  extractor,
		repository: repository,
	}
}
