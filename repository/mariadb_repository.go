package repository

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/erdnaxeli/PlayBot/types"
	_ "github.com/go-sql-driver/mysql"
)

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(dsn string) (mariaDbRepository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return mariaDbRepository{}, err
	}

	if err := db.Ping(); err != nil {
		return mariaDbRepository{}, err
	}

	return mariaDbRepository{
		db,
	}, nil
}

func (r mariaDbRepository) GetTags(musicRecordId int64) ([]string, error) {
	rows, err := r.db.Query(
		`
			select tag
			from playbot_tags
			where id = ?
		`,
		musicRecordId,
	)
	if err != nil {
		return []string{}, err
	}

	var tag string
	tags := make([]string, 0)

	for rows.Next() {
		err = rows.Scan(&tag)
		if err != nil {
			return []string{}, err
		}

		tags = append(tags, tag)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}

	return tags, nil
}

func (r mariaDbRepository) SaveMusicPost(post types.MusicPost) (int64, bool, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, false, err
	}

	recordId, isNew, err := r.insertOrUpdateMusicRecord(tx, post.MusicRecord)
	if err != nil {
		return 0, false, err
	}

	err = r.saveChannelPost(tx, recordId, post.Person, post.Channel)
	if err != nil {
		return recordId, isNew, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return 0, false, err
	}

	return recordId, isNew, nil
}

func (r mariaDbRepository) SaveTags(musicRecordId int64, tags []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		`
			insert into playbot_tags (
				id,
				tag
			)
			values (?, ?)
			on duplicate key update id=id
		`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, tag := range tags {
		_, err := stmt.Exec(musicRecordId, tag)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func (r mariaDbRepository) SearchMusicRecord(ctx context.Context, channel types.Channel, words []string, tags []string) (chan SearchResult, error) {
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
	if channel.Name != "" {
		filters = append(filters, "pc.chan = ?")
		dbArgs = append(dbArgs, channel.Name)
	}
	for _, word := range words {
		filters = append(filters, "concat(p.sender, ' ', p.title) like ?")
		dbArgs = append(dbArgs, "%"+word+"%")
	}
	for _, tag := range tags {
		filters = append(filters, "p.id in (select pt.id from playbot_tags pt where pt.tag = ?)")
		dbArgs = append(dbArgs, tag)
	}

	query += " and " + strings.Join(filters, " and ") + " "
	query += " order by rand()"

	ch := make(chan SearchResult)
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

			searchResult := SearchResult{
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

func (mariaDbRepository) insertOrUpdateMusicRecord(tx *sql.Tx, record types.MusicRecord) (int64, bool, error) {
	result, err := tx.Exec(
		`
			insert into playbot (
				type,
				url,
				sender,
				title,
				duration,
				external_id
			)
			values (
				?, ?, ?, ?, ?, ?
			)
			on duplicate key update
				id = last_insert_id(id),
				type = value(type),
				sender = value(sender),
				title = value(title),
				duration = value(duration),
				external_id = value(external_id)
		`,
		record.Source,
		record.Url,
		record.Band.Name,
		record.Name,
		int(record.Duration.Seconds()),
		record.RecordId,
	)
	if err != nil {
		return 0, false, err
	}

	recordId, err := result.LastInsertId()
	if err != nil {
		return 0, false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, false, err
	}

	return recordId, rowsAffected == 1, nil
}

func (mariaDbRepository) saveChannelPost(
	tx *sql.Tx, recordId int64, person types.Person, channel types.Channel,
) error {
	_, err := tx.Exec(
		`
			insert into playbot_chan (
				sender_irc,
				content,
				chan
			)
			values (
				?, ?, ?
			)
		`,
		person.Name,
		recordId,
		channel.Name,
	)
	if err != nil {
		return err
	}

	return nil
}
