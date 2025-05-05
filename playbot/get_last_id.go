package playbot

import "github.com/erdnaxeli/PlayBot/types"

func (p *Playbot) GetLastID(channel types.Channel, offset int) (int64, error) {
	if offset > 0 {
		return 0, ErrInvalidOffset
	}

	return p.repository.GetLastID(channel, -offset)
}
