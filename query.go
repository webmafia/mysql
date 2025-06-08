package mysql

import (
	"context"
	"database/sql"
	"unsafe"

	"github.com/webmafia/fast"
)

func (db *DB) Query(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {
	if len(args) > 0 {
		return db.queryArgs(ctx, query, fast.Noescape(args))
	}

	return db.query(ctx, query)
}

func (db *DB) queryArgs(ctx context.Context, query string, args []any) (rows *sql.Rows, err error) {
	dstQuery := db.qb.buf.Get()
	defer db.qb.buf.Put(dstQuery)

	dstArgs := db.qb.getArgs(len(args))
	defer db.qb.putArgs(dstArgs)

	if err = db.qb.buildQuery(dstQuery, dstArgs, query, args); err != nil {
		return
	}

	return db.query(ctx, dstQuery.String(), *dstArgs...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...any) (row *sql.Row) {
	if len(args) > 0 {
		return db.queryRowArgs(ctx, query, fast.Noescape(args))
	}

	return db.queryRow(ctx, query)
}

func (db *DB) queryRowArgs(ctx context.Context, query string, args []any) (row *sql.Row) {
	dstQuery := db.qb.buf.Get()
	defer db.qb.buf.Put(dstQuery)

	dstArgs := db.qb.getArgs(len(args))
	defer db.qb.putArgs(dstArgs)

	if err := db.qb.buildQuery(dstQuery, dstArgs, query, args); err != nil {
		return rowError(err)
	}

	return db.queryRow(ctx, dstQuery.String(), *dstArgs...)
}

func (db *DB) query(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {

	switch c := ctx.(type) {

	case *Tx:
		return c.query(ctx, query, args...)

	case *Conn:
		return c.conn.QueryContext(ctx, query, args...)

	default:
		return db.db.QueryContext(ctx, query, args...)
	}
}

func (db *DB) queryRow(ctx context.Context, query string, args ...any) (row *sql.Row) {

	switch c := ctx.(type) {

	case *Tx:
		return c.queryRow(ctx, query, args...)

	case *Conn:
		return c.conn.QueryRowContext(ctx, query, args...)

	default:
		return db.db.QueryRowContext(ctx, query, args...)
	}
}

func rowError(err error) *sql.Row {
	row := new(sql.Row)

	// The first field of sql.Row is an error.
	*(*error)(unsafe.Pointer(row)) = err

	return row
}
