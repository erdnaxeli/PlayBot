package playbot

import (
	"context"
	"fmt"

	"github.com/erdnaxeli/PlayBot/types"
)

// Search for a music record. It returns a channel to stream SearchResult objects.
// Closing the channel will produce a panic. If you want to notify than no more
// results will be needed, cancel the context.
func (p Playbot) SearchMusicRecord(
	ctx context.Context, channel types.Channel, words []string, tags []string,
) (chan SearchResult, error) {
	ch, err := p.repository.SearchMusicRecord(ctx, channel, words, tags)
	if err != nil {
		return nil, fmt.Errorf("error while searching music record: %w", err)
	}

	return ch, nil
}
