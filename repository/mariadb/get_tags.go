package mariadb

// GetTags returns the tags associated to a given music record.
//
// If the music record ID references a non existant music record, an empty slice is returned, and error is nil.
func (r Repository) GetTags(recordID int64) ([]string, error) {
	rows, err := r.db.Query(
		`
			select tag
			from playbot_tags
			where id = ?
		`,
		recordID,
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
