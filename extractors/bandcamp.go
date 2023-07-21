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

func NewBandcampExtractor(ldJsonExtractor ldjson.LdJsonExtractor) BandcampExtractor {
	return BandcampExtractor{
		ldJsonExtractor: ldJsonExtractor,
	}
}

func (e BandcampExtractor) Match(url string) (string, string) {
	re := regexp.MustCompile(
		`(?:^|[^!])(https?://([a-z]+)\.bandcamp.com/track/([a-zA-Z0-9_-]+))`,
	)
	groups := re.FindStringSubmatch(url)
	if groups == nil {
		return "", ""
	}

	normalizedUrl := fmt.Sprintf("https://%s.bandcamp.com/track/%s", groups[2], groups[3])
	return groups[1], normalizedUrl
}

func (e BandcampExtractor) Extract(recordId string) (types.MusicRecord, error) {
	return e.ldJsonExtractor.Extract(recordId)
}
