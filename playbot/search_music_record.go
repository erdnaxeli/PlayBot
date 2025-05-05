package playbot

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/erdnaxeli/PlayBot/types"
)

// Search represents a search request.
type Search struct {
	// Context used for the search. If the same search is execute again, the same
	// context will be used. If this context is canceled, the search is discarded
	// and a new one is started.
	Ctx context.Context

	// Channel where the search is done.
	Channel types.Channel

	// If true, the search is done globally, else it is done in the channel.
	GlobalSearch bool

	// Words to do a text based search.
	Words []string

	// Exact tags to search for.
	Tags []string

	// Tags to exclude from the search.
	ExcludedTags []string
}

// SearchMusicRecord searches for a music record.
//
// The first time a search is done in a channel, its parameters and context are
// saved. When doing another search in the same channel, if the parameters are the
// same as the previous ones, it returns the next result for the previous search.
// If the parameters are different from the previous ones, the previous search is
// discarded and a new one is started.
//
// If not result is found, an error [NoRecordFound] is returned. To distinguish between
// a search returning no result and a search where all results have been consumed you
// must look at the count returned value. If it is zero it means there is no result,
// else it means all results have been consumed.
//
// search.Ctx is the context used for the whole search. If the search is consumed again,
// the initial search.Ctx given when starting the search will be used and the current one
// will be ignored. If this context is canceled, the whole search is discarded and a new
// one is started.
// ctx is the context used to return the result. If the context is canceled no
// result is returned, but the search is kept and can be consumed again.
func (p *Playbot) SearchMusicRecord(
	ctx context.Context,
	search Search,
) (count int64, result SearchResult, err error) {
	cursor, err := p.getOrCreateSearchCursor(search)
	if err != nil {
		return 0, nil, err
	}

	result, err = p.consumeSearchCursor(ctx, search, cursor)
	return cursor.count, result, err
}

func (p *Playbot) getOrCreateSearchCursor(search Search) (searchCursor, error) {
	cursor, ok := p.getSearchCursor(search.Channel)

	if ok &&
		cursor.search.GlobalSearch == search.GlobalSearch &&
		reflect.DeepEqual(cursor.search.Words, search.Words) &&
		reflect.DeepEqual(cursor.search.Tags, search.Tags) {

		return cursor, nil
	}

	p.discardSearch(search.Channel)

	var channelToSearch types.Channel
	if !search.GlobalSearch {
		channelToSearch = search.Channel
	}

	ctx, cancel := context.WithCancel(search.Ctx)
	count, ch, err := p.repository.SearchMusicRecord(
		ctx, channelToSearch, search.Words, search.Tags,
	)
	if err != nil {
		cancel()
		return searchCursor{}, fmt.Errorf("error while searching music record: %w", err)
	}

	cursor = searchCursor{
		cancel: cancel,
		count:  count,
		ch:     ch,
		search: search,
	}

	p.setSearchCursor(search.Channel, cursor)
	return cursor, nil
}

func (p *Playbot) consumeSearchCursor(ctx context.Context, search Search, cursor searchCursor) (SearchResult, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("current search canceled: %w", ctx.Err())
	case <-cursor.search.Ctx.Done():
		p.discardSearch(search.Channel)
		return nil, SearchCanceledError{cursor.search.Ctx.Err()}
	case result, ok := <-cursor.ch:
		if !ok {
			// no more results, we discard the search
			p.discardSearch(search.Channel)

			if err := cursor.search.Ctx.Err(); err != nil {
				return nil, SearchCanceledError{err}
			}

			err := ErrNoRecordFound
			return nil, err
		}

		return result, nil
	}
}

func (p *Playbot) discardSearch(channel types.Channel) {
	cursor, ok := p.getSearchCursor(channel)
	if !ok {
		return
	}

	cursor.cancel()
	p.deleteSearchCursor(channel)
}

func (p *Playbot) getSearchCursor(channel types.Channel) (searchCursor, bool) {
	p.searchesMutex.RLock()
	defer p.searchesMutex.RUnlock()

	cursor, ok := p.searches[channel]
	return cursor, ok
}

func (p *Playbot) setSearchCursor(channel types.Channel, cursor searchCursor) {
	p.searchesMutex.Lock()
	defer p.searchesMutex.Unlock()

	p.searches[channel] = cursor
	context.AfterFunc(cursor.search.Ctx, func() {
		log.Printf("Discard canceled search for channel %s", channel)
		p.deleteSearchCursor(channel)
	})
}

func (p *Playbot) deleteSearchCursor(channel types.Channel) {
	p.searchesMutex.Lock()
	defer p.searchesMutex.Unlock()

	delete(p.searches, channel)
}
