package extractors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/types"
)

type SoundCloudExtractor struct {
	ldJsonExtractor ldjson.LdJsonExtractor
}

func NewSoundCloudExtractor(ldJsonExtractor ldjson.LdJsonExtractor) SoundCloudExtractor {
	return SoundCloudExtractor{
		ldJsonExtractor: ldJsonExtractor,
	}
}

func (SoundCloudExtractor) Match(url string) (string, string) {
	re := regexp.MustCompile(`(?:^|[^!])(https?://(?:www\.)?soundcloud.com/([a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+)(?:\?.+)?)`)
	groups := re.FindStringSubmatch(url)
	if groups == nil {
		return "", ""
	}

	return groups[1], groups[2]
}

func (e SoundCloudExtractor) Extract(recordId string) (types.MusicRecord, error) {
	record, err := e.ldJsonExtractor.Extract("https://m.soundcloud.com/" + recordId)
	if err != nil {
		return types.MusicRecord{}, err
	}

	recordId, err = e.getRecordId(record.RecordId)
	if err != nil {
		return types.MusicRecord{}, err
	}

	record.RecordId = recordId
	record.Source = "soundcloud"
	return record, nil
}

func (SoundCloudExtractor) getRecordId(recordId string) (string, error) {
	parts := strings.Split(recordId, ":")
	if len(parts) != 3 {
		return "", fmt.Errorf("unknown Soundcloud record id '%s'", recordId)
	}

	return parts[2], nil
}
