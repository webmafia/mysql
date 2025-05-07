package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/webmafia/fast/buffer"
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
	}

	return
}

var _ context.Context = (*Tx)(nil)

type Tx struct {
	context.Context
	db     *DB
	tx     *sql.Tx
	sp     savepoint
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

	if tx.sp > 0 {
		return tx.releaseSavepoint()
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
	tx.Context = nil
	tx.db = nil
	tx.tx = nil
	tx.sp = 0
	tx.closed = false

	tx.db.txPool.Put(tx)
}
