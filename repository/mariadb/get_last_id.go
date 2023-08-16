package mariadb

import (
	"database/sql"
	"errors"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/types"
)

func (r mariaDbRepository) GetLastID(channel types.Channel, offset int) (int64, error) {
	row := r.db.QueryRow(
		`
			select p.id
			from playbot p
			join playbot_chan pc
				on p.id = pc.content
			where
				pc.chan = ?
			order by pc.date desc, pc.id desc
			limit 1
			offset ?
		`,
		channel.Name,
		offset,
	)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, playbot.NoRecordFoundError
		}
		return 0, err
	}

	return id, err
}
