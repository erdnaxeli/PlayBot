package mariadb

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
