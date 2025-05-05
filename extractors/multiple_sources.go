package extractors

import "github.com/erdnaxeli/PlayBot/types"

// MultipleSourcesExtractor allows to test multiple extractors.
type MultipleSourcesExtractor struct {
	extractors []MatcherExtractor
}

// New create a new instance of a MultipleSourcesExtractor object.
func New(extractors ...MatcherExtractor) MultipleSourcesExtractor {
	return MultipleSourcesExtractor{
		extractors: extractors,
	}
}

// Extract iterates over all sources and use the first one matching.
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
