package mariadb

import (
	"database/sql"
	"errors"
	"time"

	"github.com/erdnaxeli/PlayBot/types"
)

func (r mariaDbRepository) GetMusicRecord(musicRecordId int64) (types.MusicRecord, error) {
	row := r.db.QueryRow(
		`
			select
				p.sender,
				p.title,
				p.url,
				p.duration,
				p.external_id,
				p.type
			from playbot p
			where p.id = ?
		`,
		musicRecordId,
	)

	var title, url, source string
	var recordID, sender sql.Null[string]
	var duration sql.Null[int64]
	err := row.Scan(&sender, &title, &url, &duration, &recordID, &source)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.MusicRecord{}, nil
		}

		return types.MusicRecord{}, err
	}

	record := types.MusicRecord{
		Band:     types.Band{Name: sender.V},
		Duration: time.Duration(duration.V * int64(time.Second)),
		Name:     title,
		RecordID: recordID.V,
		Source:   source,
		URL:      url,
	}
	return record, nil
}
