package extractors

import "github.com/erdnaxeli/PlayBot/types"

var extractors []Extractor

func Extract(url string) (string, types.MusicRecord, error) {
	for _, extractor := range extractors {
		if matchedUrl, recordId := extractor.Match(url); recordId != "" {
			record, err := extractor.Extract(recordId)
			if err != nil {
				return "", types.MusicRecord{}, err
			} else {
				return matchedUrl, record, nil
			}
		}
	}

	return "", types.MusicRecord{}, &UnknownRecordSourceError{}
}
