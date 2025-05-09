package extractors

import (
	"fmt"
	"regexp"

	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/types"
)

// BandcampExtractor implements the MatcherExtractor interface.
type BandcampExtractor struct {
	ldJSONExtractor ldjson.Extractor
}

// NewBandcampExtractor return a new BandcampExtractor instance.
func NewBandcampExtractor(ldJSONExtractor ldjson.Extractor) BandcampExtractor {
	return BandcampExtractor{
		ldJSONExtractor: ldJSONExtractor,
	}
}

// Match returns the URL matched and the record ID.
func (e BandcampExtractor) Match(url string) (string, string) {
	re := regexp.MustCompile(
		`(?:^|[^!])(https?://([a-z]+)\.bandcamp\.com/track/([a-zA-Z0-9_-]+))`,
	)
	groups := re.FindStringSubmatch(url)
	if groups == nil {
		return "", ""
	}

	normalizedURL := fmt.Sprintf("https://%s.bandcamp.com/track/%s", groups[2], groups[3])
	return groups[1], normalizedURL
}

// Extract return a record data.
func (e BandcampExtractor) Extract(recordID string) (types.MusicRecord, error) {
	record, err := e.ldJSONExtractor.Extract(recordID)
	if err != nil {
		return types.MusicRecord{}, err
	}

	record.Source = "bandcamp"
	return record, nil
}
