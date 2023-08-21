package mariadb

import (
	"errors"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/go-sql-driver/mysql"
)

const (
	_ER_DUP_ENTRY           = 1062
	_ER_NO_REFERENCED_ROW   = 1216
	_ER_NO_REFERENCED_ROW_2 = 1452
)

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
		mysqlErr := &mysql.MySQLError{}
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == _ER_DUP_ENTRY {
				return nil
			} else if mysqlErr.Number == _ER_NO_REFERENCED_ROW ||
				mysqlErr.Number == _ER_NO_REFERENCED_ROW_2 {

				return playbot.NoRecordFoundError
			}
		}

		return err
	}

	return nil
}
