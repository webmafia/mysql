package mysql

import (
	"database/sql"
	"sync"

	"github.com/go-sql-driver/mysql"
)

type DB struct {
	db      *sql.DB
	qb      queryBuilder
	txPool  sync.Pool
	valPool sync.Pool
	lruPool sync.Pool
}

// var defaultLogger = mysql.Logger(log.New(os.Stderr, "[mysql] ", log.Ldate|log.Ltime))

type Config = mysql.Config

func NewConfig() *Config {
	return mysql.NewConfig()
}

func New(cfg *Config) (db *DB, err error) {
	db = &DB{}
	c, err := mysql.NewConnector(cfg)

	if err != nil {
		return
	}

	// if cfg.Loc == nil {
	// 	cfg.Loc = time.UTC
	// }

	// if cfg.MaxAllowedPacket <= 0 {
	// 	cfg.MaxAllowedPacket = 64 << 20 // 64 MiB
	// }

	// if cfg.Logger == nil {
	// 	cfg.Logger = defaultLogger
	// }

	// cfg.AllowNativePasswords = true
	// cfg.CheckConnLiveness = true

	db.db = sql.OpenDB(c)
	return
}

func NewFromString(connStr string) (db *DB, err error) {
	cfg, err := mysql.ParseDSN(connStr)

	if err != nil {
		return
	}

	return New(cfg)
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Stats() sql.DBStats {
	return db.db.Stats()
}
