package extractors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/types"
)

// SoundCloudUnknownRecordIDError is the error when a record ID format is unknown.
type SoundCloudUnknownRecordIDError struct {
	// RecordID is the faulty record ID.
	RecordID string
}

func (err SoundCloudUnknownRecordIDError) Error() string {
	return fmt.Sprintf("unknown SoundCloud record ID '%s'", err.RecordID)
}

// SoundCloudExtractor implements the MatcherExtractor interface.
type SoundCloudExtractor struct {
	ldJSONExtractor ldjson.Extractor
}

// NewSoundCloudExtractor returns a new instance of SoundCloudExtractor.
func NewSoundCloudExtractor(ldJSONExtractor ldjson.Extractor) SoundCloudExtractor {
	return SoundCloudExtractor{
		ldJSONExtractor: ldJSONExtractor,
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

// Extract returns the record data.
func (e SoundCloudExtractor) Extract(recordID string) (types.MusicRecord, error) {
	record, err := e.ldJSONExtractor.Extract("https://m.soundcloud.com/" + recordID)
	if err != nil {
		return types.MusicRecord{}, err
	}

	recordID, err = e.getRecordID(record.RecordID)
	if err != nil {
		return types.MusicRecord{}, err
	}

	record.RecordID = recordID
	record.Source = "soundcloud"
	return record, nil
}

func (SoundCloudExtractor) getRecordID(recordID string) (string, error) {
	parts := strings.Split(recordID, ":")
	if len(parts) != 3 {
		return "", SoundCloudUnknownRecordIDError{recordID}
	}

	return parts[2], nil
}
