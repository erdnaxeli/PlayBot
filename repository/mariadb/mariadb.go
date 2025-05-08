// Package mariadb implements a playbot.Repository using a MariaDB server.
package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/erdnaxeli/PlayBot/types"
	// register mysql driver
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

// New returns a new instance of a repository.
//
// It connect to the db using the given parameters.
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

// NewFromDB returns a new instance of a repository.
//
// It takes an *sql.DB instance as a parameter.
func NewFromDB(db *sql.DB) (mariaDbRepository, error) {
	if err := db.Ping(); err != nil {
		return mariaDbRepository{}, err
	}

	return mariaDbRepository{
		db,
	}, nil
}
