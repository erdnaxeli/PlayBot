package mariadb

import (
	"database/sql"
	"errors"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/types"
)

// GetLastID returns the id of the last post known for a channel.
//
// The offset can be used to select a previous post.
// Offset 0 means the last post, 1 the second last post, and so on.
func (r Repository) GetLastID(channel types.Channel, offset int) (int64, error) {
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
			return 0, playbot.ErrNoRecordFound
		}
		return 0, err
	}

	return id, err
}
