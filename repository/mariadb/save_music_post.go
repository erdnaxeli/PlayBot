package mariadb

import (
	"database/sql"

	"github.com/erdnaxeli/PlayBot/types"
)

func (r mariaDbRepository) SaveMusicPost(post types.MusicPost) (int64, bool, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, false, err
	}

	recordID, isNew, err := r.insertOrUpdateMusicRecord(tx, post.MusicRecord)
	if err != nil {
		return 0, false, err
	}

	err = r.saveChannelPost(tx, recordID, post.Person, post.Channel)
	if err != nil {
		return recordID, isNew, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return 0, false, err
	}

	return recordID, isNew, nil
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
		record.URL,
		record.Band.Name,
		record.Name,
		int(record.Duration.Seconds()),
		record.RecordID,
	)
	if err != nil {
		return 0, false, err
	}

	recordID, err := result.LastInsertId()
	if err != nil {
		return 0, false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, false, err
	}

	return recordID, rowsAffected == 1, nil
}

func (mariaDbRepository) saveChannelPost(
	tx *sql.Tx, recordID int64, person types.Person, channel types.Channel,
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
		recordID,
		channel.Name,
	)
	if err != nil {
		return err
	}

	return nil
}
