package playbot

import "fmt"

// GetTags returns the tag of a music record.
func (p *Playbot) GetTags(musicRecordID int64) ([]string, error) {
	tags, err := p.repository.GetTags(musicRecordID)
	if err != nil {
		return tags, fmt.Errorf("error while retrieving tags: %w", err)
	}

	return tags, nil
}
