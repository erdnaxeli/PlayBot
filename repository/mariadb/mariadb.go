package mariadb

import (
	"database/sql"
	"fmt"

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

func (s searchResult) ID() int64 {
	return s.id
}

func (s searchResult) MusicRecord() types.MusicRecord {
	return s.musicRecord
}

func New(user string, password string, host string, dbname string) (mariaDbRepository, error) {
	dsn := fmt.Sprintf(
		"%s:%s@(%s)/%s?parseTime=true&loc=Europe%%2FParis",
		user, password, host, dbname,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return mariaDbRepository{}, err
	}

	return NewFromDB(db)
}

func NewFromDB(db *sql.DB) (mariaDbRepository, error) {
	if err := db.Ping(); err != nil {
		return mariaDbRepository{}, err
	}

	return mariaDbRepository{
		db,
	}, nil
}
