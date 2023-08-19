package mariadb

func (r mariaDbRepository) SaveFav(user string, recordID int64) error {
	_, err := r.db.Exec(
		`
			insert into playbot_fav (user, id)
			value (?, ?)
		`,
		user,
		recordID,
	)
	if err != nil {
		return err
	}

	return nil
}
