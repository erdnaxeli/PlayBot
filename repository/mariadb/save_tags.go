package mariadb

// SaveTags add some tags to a music record.
//
// The given tags are added to any existing tags already linked to this music record.
func (r Repository) SaveTags(musicRecordID int64, tags []string) error {
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
	defer func() { _ = stmt.Close() }()

	for _, tag := range tags {
		_, err := stmt.Exec(musicRecordID, tag)
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
