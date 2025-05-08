package mariadb

import (
	"errors"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/go-sql-driver/mysql"
)

const (
	dupEntryErr         = 1062
	noReferencedRowErr  = 1216
	noReferencedRow2Err = 1452
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
			switch mysqlErr.Number {
			case dupEntryErr:
				return nil
			case noReferencedRowErr, noReferencedRow2Err:
				return playbot.ErrNoRecordFound
			}
		}

		return err
	}

	return nil
}
