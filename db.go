package mysql

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db     *sql.DB
	qb     queryBuilder
	txPool sync.Pool
}

func New(connStr string) (db *DB, err error) {
	db = &DB{}

	if db.db, err = sql.Open("mysql", connStr); err != nil {
		return
	}

	return
}

func (db *DB) Close() error {
	return db.db.Close()
}
