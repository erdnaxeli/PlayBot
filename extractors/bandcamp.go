package extractors

import (
	"fmt"
	"regexp"

	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/types"
)

type BandcampExtractor struct {
	ldJsonExtractor ldjson.LdJsonExtractor
}

// NewBandcampExtractor return a new BandcampExtractor instance.
func NewBandcampExtractor(ldJsonExtractor ldjson.LdJsonExtractor) BandcampExtractor {
	return BandcampExtractor{
		ldJsonExtractor: ldJsonExtractor,
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

	normalizedUrl := fmt.Sprintf("https://%s.bandcamp.com/track/%s", groups[2], groups[3])
	return groups[1], normalizedUrl
}

// Extract return a record data.
func (e BandcampExtractor) Extract(recordId string) (types.MusicRecord, error) {
	record, err := e.ldJsonExtractor.Extract(recordId)
	if err != nil {
		return types.MusicRecord{}, err
	}

	record.Source = "bandcamp"
	return record, nil
}
