package mariadb

import (
	"database/sql"

	"github.com/erdnaxeli/PlayBot/types"
	_ "github.com/go-sql-driver/mysql"
)

type mariaDbRepository struct {
	db *sql.DB
}

type searchResult struct {
	id          int64
	musicRecord types.MusicRecord
}

func (s searchResult) Id() int64 {
	return s.id
}

func (s searchResult) MusicRecord() types.MusicRecord {
	return s.musicRecord
}

func New(dsn string) (mariaDbRepository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return mariaDbRepository{}, err
	}

	if err := db.Ping(); err != nil {
		return mariaDbRepository{}, err
	}

	return mariaDbRepository{
		db,
	}, nil
}
