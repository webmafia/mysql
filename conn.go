package mysql

import (
	"context"
	"database/sql"
)

func (db *DB) Conn(ctx context.Context) (c *Conn, err error) {

	// Already a connection
	if conn, ok := ctx.(*Conn); ok {
		return conn, nil
	}

	conn, err := db.db.Conn(ctx)

	if err != nil {
		return
	}

	return &Conn{
		Context: ctx,
		conn:    conn,
	}, nil
}

var _ context.Context = (*Conn)(nil)

type Conn struct {
	context.Context
	conn *sql.Conn
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
