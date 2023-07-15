package extractors

import (
	"github.com/erdnaxeli/PlayBot/types"
)

type Extractor interface {
	// Match the given url to the format expected by this Extractor. If it matches
	// it returns a tuple with the whole URL matched and the unique identifier for this
	// music recourd, else it returns two empty string.
	Match(url string) (string, string)

	// Extract the record data for the given music record id.
	Extract(musicRecordId string) (types.MusicRecord, error)
}
