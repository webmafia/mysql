package mysql

import (
	"context"
	"database/sql"

	"github.com/webmafia/fast"
)

func (db *DB) Exec(ctx context.Context, query string, args ...any) (result sql.Result, err error) {
	if len(args) > 0 {
		return db.execArgs(ctx, query, fast.Noescape(args))
	}

	switch c := ctx.(type) {

	case *Tx:
		return c.exec(ctx, query, args...)

	case *Conn:
		return c.conn.ExecContext(ctx, query, args...)

	default:
		return db.db.ExecContext(ctx, query, args...)
	}
}

func (db *DB) execArgs(ctx context.Context, query string, args []any) (result sql.Result, err error) {
	dstQuery := db.qb.buf.Get()
	defer db.qb.buf.Put(dstQuery)

	dstArgs := db.qb.getArgs(len(args))
	defer db.qb.putArgs(dstArgs)

	if err = db.qb.buildQuery(dstQuery, dstArgs, query, args); err != nil {
		return
	}

	switch c := ctx.(type) {

	case *Tx:
		return c.exec(ctx, dstQuery.String(), *dstArgs...)

	case *Conn:
		return c.conn.ExecContext(ctx, dstQuery.String(), *dstArgs...)
	}

	return db.db.ExecContext(ctx, dstQuery.String(), *dstArgs...)
}
