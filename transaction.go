package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/webmafia/fast/buffer"
	"github.com/webmafia/lru"
)

var ErrReleasedTransaction = errors.New("tried to operate on a released transaction")

func (db *DB) Transaction(ctx context.Context, readOnly ...bool) (tx *Tx, err error) {
	tx, ok := db.txPool.Get().(*Tx)

	if !ok {
		tx = new(Tx)
	}

	tx.Context = ctx
	tx.db = db

	if parent, ok := ctx.(*Tx); ok {
		tx.child = true
		tx.stmts = parent.stmts
		tx.tx = parent.tx
		tx.sp = parent.sp + 1

		if err = tx.savepoint(); err != nil {
			return
		}
	} else {
		var opts *sql.TxOptions

		if len(readOnly) > 0 && readOnly[0] {
			opts = &sql.TxOptions{ReadOnly: true}
		}

		if tx.tx, err = db.db.BeginTx(ctx, opts); err != nil {
			return
		}

		stmts, ok := db.lruPool.Get().(lru.LRU[uint64, *sql.Stmt])

		if !ok {
			stmts = lru.NewThreadSafe(64, func(_ uint64, stmt *sql.Stmt) {
				if err := stmt.Close(); err != nil {
					log.Println(err)
				}
			})
		}

		tx.stmts = stmts
	}

	return
}

var _ context.Context = (*Tx)(nil)

type Tx struct {
	context.Context
	db     *DB
	tx     *sql.Tx
	stmts  lru.LRU[uint64, *sql.Stmt]
	sp     savepoint
	child  bool
	closed bool
}

var _ StringEncoder = savepoint(0)

type savepoint int

// EncodeString implements StringEncoder.
func (s savepoint) EncodeString(b *buffer.Buffer) {
	b.WriteString("tx_sp_")
	b.Str().WriteInt(int(s))
}

func (tx *Tx) savepoint() (err error) {
	_, err = tx.db.Exec(tx, "SAVEPOINT %T", tx.sp)
	return
}

func (tx *Tx) rollbackSavepoint() (err error) {
	_, err = tx.db.Exec(tx, "ROLLBACK TO SAVEPOINT %T", tx.sp)
	return
}

func (tx *Tx) releaseSavepoint() (err error) {
	_, err = tx.db.Exec(tx, "RELEASE SAVEPOINT %T", tx.sp)
	return
}

func (tx *Tx) Commit() (err error) {
	if tx.closed {
		return ErrReleasedTransaction
	}

	if err = tx.tx.Commit(); err != nil {
		return
	}

	tx.closed = true
	return
}

func (tx *Tx) Release() (err error) {
	defer tx.release()

	if tx.closed {
		return
	}

	if tx.sp > 0 {
		return tx.rollbackSavepoint()
	}

	if err = tx.tx.Rollback(); err != nil {
		return
	}

	tx.closed = true
	return
}

func (tx *Tx) release() {
	db := tx.db

	if !tx.child {
		tx.stmts.Reset()
		tx.db.lruPool.Put(tx.stmts)
	} else {
		tx.child = false
	}

	tx.Context = nil
	tx.db = nil
	tx.tx = nil
	tx.sp = 0
	tx.closed = false
	tx.stmts = nil

	db.txPool.Put(tx)
}

func (tx *Tx) exec(ctx context.Context, query string, args ...any) (_ sql.Result, err error) {
	if len(args) == 0 {
		return tx.tx.ExecContext(ctx, query)
	}

	stmt, err := tx.stmt(ctx, query)

	if err != nil {
		return
	}

	return stmt.ExecContext(ctx, args...)
}

func (tx *Tx) query(ctx context.Context, query string, args ...any) (_ *sql.Rows, err error) {
	if len(args) == 0 {
		return tx.tx.QueryContext(ctx, query)
	}

	stmt, err := tx.stmt(ctx, query)

	if err != nil {
		return
	}

	return stmt.QueryContext(ctx, args...)
}

func (tx *Tx) queryRow(ctx context.Context, query string, args ...any) *sql.Row {
	if len(args) == 0 {
		return tx.tx.QueryRowContext(ctx, query)
	}

	stmt, err := tx.stmt(ctx, query)

	if err != nil {
		return rowError(err)
	}

	return stmt.QueryRowContext(ctx, args...)
}

func (tx *Tx) stmt(ctx context.Context, query string) (*sql.Stmt, error) {
	return tx.stmts.GetOrSet(xxhash.Sum64String(query), func(key uint64) (*sql.Stmt, error) {
		return tx.tx.PrepareContext(ctx, strings.Clone(query))
	})
}
