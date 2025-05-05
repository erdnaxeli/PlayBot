package extractors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/types"
)

type SoundCloudUnknownRecordIDError struct {
	RecordID string
}

func (err SoundCloudUnknownRecordIDError) Error() string {
	return fmt.Sprintf("unknown SoundCloud record ID '%s'", err.RecordID)
}

type SoundCloudExtractor struct {
	ldJsonExtractor ldjson.LdJsonExtractor
}

func NewSoundCloudExtractor(ldJsonExtractor ldjson.LdJsonExtractor) SoundCloudExtractor {
	return SoundCloudExtractor{
		ldJsonExtractor: ldJsonExtractor,
	}
}

// Match returns the URL matched and the record ID.
func (SoundCloudExtractor) Match(url string) (string, string) {
	re := regexp.MustCompile(`(?:^|[^!])(https?://(?:www\.)?soundcloud\.com/([a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+)(?:\?.+)?)`)
	groups := re.FindStringSubmatch(url)
	if groups == nil {
		return "", ""
	}

	return groups[1], groups[2]
}

// Extracts returns the record data.
func (e SoundCloudExtractor) Extract(recordId string) (types.MusicRecord, error) {
	record, err := e.ldJsonExtractor.Extract("https://m.soundcloud.com/" + recordId)
	if err != nil {
		return types.MusicRecord{}, err
	}

	recordId, err = e.getRecordId(record.RecordID)
	if err != nil {
		return types.MusicRecord{}, err
	}

	record.RecordID = recordId
	record.Source = "soundcloud"
	return record, nil
}

func (SoundCloudExtractor) getRecordId(recordId string) (string, error) {
	parts := strings.Split(recordId, ":")
	if len(parts) != 3 {
		return "", SoundCloudUnknownRecordIDError{recordId}
	}

	return parts[2], nil
}
