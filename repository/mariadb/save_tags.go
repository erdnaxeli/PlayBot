package mariadb

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
