// Package extractors implements multiple types to extract music record data from different sources.
package extractors

import (
	"github.com/erdnaxeli/PlayBot/types"
)

// MatcherExtractor is the interface that must be implemented by any extractor.
//
// An extractor is an object offering method to extract the record data from an URL to a music record.
type MatcherExtractor interface {
	// Match the given url to the format expected by this Extractor. If it matches
	// it returns a tuple with the whole URL matched and the unique identifier for this
	// music record, else it returns two empty string.
	Match(url string) (string, string)

	// Extract the record data for the given music record id.
	Extract(musicRecordID string) (types.MusicRecord, error)
}
