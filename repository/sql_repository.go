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

func (r sqlRepository) SaveMusicPost(post types.MusicPost) (int64, error) {
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
		_ = tx.Rollback()
		return 0, err
	}

	return recordId, nil
}

func (r sqlRepository) SaveTags(musicRecordId int64, tags []string) error {
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
