package mariadb

import (
	"context"
	"database/sql"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/types"
)

func (r mariaDbRepository) SearchMusicRecord(
	ctx context.Context, channel types.Channel, words []string, tags []string,
) (int64, chan playbot.SearchResult, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, nil, err
	}

	// We can't use defer because SearchMusicRecord will end but the goroutine still
	// needs the transaction to read rows from it.
	runtime.SetFinalizer(tx, func(tx *sql.Tx) { _ = tx.Rollback() })

	queryCount, dbArgs := makeSearchQuery(true, channel.Name, words, tags)
	row := tx.QueryRow(queryCount, dbArgs...)
	var count int64
	err = row.Scan(&count)
	if err != nil {
		return 0, nil, err
	}

	query, dbArgs := makeSearchQuery(false, channel.Name, words, tags)
	ch := make(chan playbot.SearchResult)
	rows, err := tx.Query(query, dbArgs...)
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
				log.Print("search canceled")
				return
			case ch <- searchResult:
			}
		}

		// We need the transaction to stay alive until here to be able to read results
		// from it.
		runtime.KeepAlive(tx)
		log.Print("search done")
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
