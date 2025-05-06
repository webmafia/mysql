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

	return db.db.ExecContext(ctx, query)
}

func (db *DB) execArgs(ctx context.Context, query string, args []any) (result sql.Result, err error) {
	dstQuery := db.qb.buf.Get()
	defer db.qb.buf.Put(dstQuery)

	dstArgs := db.qb.getArgs(len(args))
	defer db.qb.putArgs(dstArgs)

	if err = db.qb.buildQuery(dstQuery, dstArgs, query, args); err != nil {
		return
	}

	return db.db.ExecContext(ctx, dstQuery.String(), *dstArgs...)
}
