package playbot

import "fmt"

// SaveTags add some tags to a music record
func (p *Playbot) SaveTags(musicRecordID int64, tags []string) error {
	err := p.repository.SaveTags(musicRecordID, tags)
	if err != nil {
		return fmt.Errorf("error while saving tags: %w", err)
	}

	return nil
}
