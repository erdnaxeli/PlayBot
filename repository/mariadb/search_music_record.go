package mariadb

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/types"
)

func (r mariaDbRepository) SearchMusicRecord(
	ctx context.Context, channel types.Channel, words []string, tags []string,
) (int64, chan playbot.SearchResult, error) {
	queryCount, dbArgs := makeSearchQuery(true, channel.Name, words, tags)
	row := r.db.QueryRow(queryCount, dbArgs...)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, nil, err
	}

	query, dbArgs := makeSearchQuery(false, channel.Name, words, tags)
	ch := make(chan playbot.SearchResult)
	rows, err := r.db.Query(query, dbArgs...)
	if err != nil {
		return 0, nil, err
	}

	go func() {
		defer func() { _ = rows.Close() }()
		defer close(ch)

		for rows.Next() {
			var id int64
			var title, url, source string
			var recordID, sender sql.NullString
			var duration sql.NullInt64
			err := rows.Scan(&id, &sender, &title, &url, &duration, &recordID, &source)
			if err != nil {
				log.Printf("Error while fetching rows after search: %s", err)
				return
			}

			searchResult := searchResult{
				id,
				types.MusicRecord{
					Band:     types.Band{Name: sender.String},
					Duration: time.Duration(duration.Int64 * int64(time.Second)),
					Name:     title,
					RecordId: recordID.String,
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

	return count, ch, nil
}

func makeSearchQuery(count bool, channelName string, words []string, tags []string) (string, []any) {
	var query string
	if count {
		query = "select count(distinct p.id)"
	} else {
		query = `
			select distinct
				p.id,
				p.sender,
				p.title,
				p.url,
				p.duration,
				p.external_id,
				p.type
		`
	}

	query += `
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

	if len(filters) > 0 {
		query += " and " + strings.Join(filters, " and ") + " "
	}

	query += " order by rand()"

	return query, dbArgs
}
