package repository

import (
	"database/sql"

	"github.com/erdnaxeli/PlayBot/types"
	_ "github.com/go-sql-driver/mysql"
)

type sqlRepository struct {
	db *sql.DB
}

func NewSqlRepository(dsn string) (sqlRepository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return sqlRepository{}, err
	}

	if err := db.Ping(); err != nil {
		return sqlRepository{}, err
	}

	return sqlRepository{
		db,
	}, nil
}

func (r sqlRepository) SaveMusicRecord(post types.MusicPost) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	recordId, err := r.insertOrUpdateMusicRecord(tx, post.MusicRecord)
	if err != nil {
		return 0, err
	}

	err = r.saveChannelPost(tx, recordId, post.Person, post.Channel)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return recordId, nil
}

func SaveTags(int, []string) error {
	return nil
}

func (sqlRepository) insertOrUpdateMusicRecord(tx *sql.Tx, record types.MusicRecord) (int64, error) {
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
		return 0, err
	}

	recordId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return recordId, nil
}

func (sqlRepository) saveChannelPost(
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
