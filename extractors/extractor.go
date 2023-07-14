package extractors

import (
	"github.com/erdnaxeli/PlayBot/types"
)

type Extractor interface {
	// Match the given url to the format expected by this Extractor. If it matches
	// return a unique identifier for this music recourd, else return an empty string.
	Match(url string) string

	// Extract the record data for the given music record id.
	Extract(musicRecordId string) (types.MusicRecord, error)
}
