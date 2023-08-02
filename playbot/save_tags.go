package playbot

import "fmt"

func (p *Playbot) SaveTags(musicRecordId int64, tags []string) error {
	err := p.repository.SaveTags(musicRecordId, tags)
	if err != nil {
		return fmt.Errorf("error while saving tags: %w", err)
	}

	return nil
}
