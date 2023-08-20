package extractors

import "github.com/erdnaxeli/PlayBot/types"

type MultipleSourcesExtractor struct {
	extractors []MatcherExtractor
}

func New(extractors ...MatcherExtractor) MultipleSourcesExtractor {
	return MultipleSourcesExtractor{
		extractors: extractors,
	}
}

func (e MultipleSourcesExtractor) Extract(url string) (types.MusicRecord, error) {
	for _, extractor := range e.extractors {
		if _, recordId := extractor.Match(url); recordId != "" {
			record, err := extractor.Extract(recordId)
			if err != nil {
				return types.MusicRecord{}, err
			} else {
				return record, nil
			}
		}
	}

	return types.MusicRecord{}, nil
}
