package mariadb

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/erdnaxeli/PlayBot/types"
)

func (r mariaDbRepository) SearchMusicRecord(
	ctx context.Context, channel types.Channel, words []string, tags []string,
) (chan searchResult, error) {
	query, dbArgs := makeSearchQuery(channel.Name, words, tags)
	ch := make(chan searchResult)
	rows, err := r.db.Query(query, dbArgs...)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() { _ = rows.Close() }()
		defer close(ch)

		for rows.Next() {
			var id, duration int64
			var sender, title, url, recordId, source string
			err := rows.Scan(&id, &sender, &title, &url, &duration, &recordId, &source)
			if err != nil {
				log.Printf("Error while fetching rows after search: %s", err)
				return
			}

			searchResult := searchResult{
				id,
				types.MusicRecord{
					Band:     types.Band{Name: sender},
					Duration: time.Duration(duration * int64(time.Second)),
					Name:     title,
					RecordId: recordId,
					Source:   source,
					Url:      url,
				},
			}

			select {
			case <-ctx.Done():
				return
			case ch <- searchResult:
			}
		}
	}()

	return ch, nil

}

func makeSearchQuery(channelName string, words []string, tags []string) (string, []any) {
	query := `
		select distinct
			p.id,
			p.sender,
			p.title,
			p.url,
			p.duration,
			p.external_id,
			p.type
		from playbot p
		join playbot_chan pc
			on p.id = pc.content
		where
			p.playlist is false
	`

	var dbArgs []any
	var filters []string
	if channelName != "" {
		filters = append(filters, "pc.chan = ?")
		dbArgs = append(dbArgs, channelName)
	}
	for _, word := range words {
		filters = append(filters, "concat(p.sender, ' ', p.title) like ?")
		dbArgs = append(dbArgs, "%"+word+"%")
	}
	for _, tag := range tags {
		filters = append(
			filters,
			"p.id in (select pt.id from playbot_tags pt where pt.tag = ?)",
		)
		dbArgs = append(dbArgs, tag)
	}

	query += " and " + strings.Join(filters, " and ") + " "
	query += " order by rand()"

	return query, dbArgs
}
