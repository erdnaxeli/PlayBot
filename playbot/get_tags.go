package playbot

import "fmt"

func (p Playbot) GetTags(musicRecordId int64) ([]string, error) {
	tags, err := p.repository.GetTags(musicRecordId)
	if err != nil {
		return tags, fmt.Errorf("error while retrieving tags: %w", err)
	}

	return tags, nil
}
