package mariadb

import (
	"database/sql"
	"errors"
)

var UserNotFoundErr = errors.New("user not found")

func (r mariaDbRepository) GetUserFromNick(nick string) (string, error) {
	row := r.db.QueryRow(
		`
			select user
			from playbot_codes
			where nick = ?
		`,
		nick,
	)

	var user string
	err := row.Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	return user, nil
}

func (r mariaDbRepository) GetUserFromCode(code string) (string, error) {
	row := r.db.QueryRow(
		`
			select user
			from playbot_codes
			where code = ?
		`,
		code,
	)

	var user string
	err := row.Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	return user, nil
}

func (r mariaDbRepository) SaveAssociation(user string, nick string) error {
	result, err := r.db.Exec(
		`
			update playbot_codes
			set
				nick = ?,
				date = now()
			where user = ?
		`,
		nick,
		user,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return UserNotFoundErr
	}

	return nil
}
