package mariadb

import (
	"database/sql"

	"github.com/erdnaxeli/PlayBot/types"
)

// SaveMusicPost saves a post.
//
// It creates or update the music record, and then create a new post.
// It returns the music record ID, a bool indicating if the music record was created or not, and an error if any.
//
// If there is any error, nothing is saved to the database, neither the music record (created or updated) nor the post.
func (r Repository) SaveMusicPost(post types.MusicPost) (int64, bool, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, false, err
	}
	defer func() { _ = tx.Rollback() }()

	recordID, isNew, err := r.insertOrUpdateMusicRecord(tx, post.MusicRecord)
	if err != nil {
		return 0, false, err
	}

	err = r.saveChannelPost(tx, recordID, post.Person, post.Channel)
	if err != nil {
		return 0, false, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, false, err
	}

	return recordID, isNew, nil
}

func (Repository) insertOrUpdateMusicRecord(tx *sql.Tx, record types.MusicRecord) (int64, bool, error) {
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

func (Repository) saveChannelPost(
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
	return err
}
