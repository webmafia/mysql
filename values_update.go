package mysql

import (
	"context"
)

func (db *DB) UpdateValues(ctx context.Context, table StringEncoder, vals *Values, cond QueryEncoder) (count int64, err error) {
	if vals.Empty() {
		return
	}

	cmd, err := db.Exec(ctx,
		"UPDATE %T SET %T WHERE %T",
		table,
		vals,
		cond,
	)

	if err == nil {
		vals.reset()
		count, err = cmd.RowsAffected()
	}

	return
}
