package mysql

import (
	"context"
	"slices"

	"github.com/webmafia/fast/buffer"
)

type InsertOptions struct {
	OnConflict func(vals *Values) EncodeQuery
	Ignore     bool // Ignore any errors.
}

func (db *DB) InsertValues(ctx context.Context, table StringEncoder, vals *Values, options ...InsertOptions) (count int64, err error) {
	if vals.Empty() {
		return
	}

	var onConflict QueryEncoder = queryEncoderNoop
	var ignore bool

	if len(options) > 0 {
		if options[0].OnConflict != nil {
			onConflict = options[0].OnConflict(vals)
		}

		ignore = options[0].Ignore
	}

	var query string

	if ignore {
		query = "INSERT IGNORE INTO %T SET %T %T"
	} else {
		query = "INSERT INTO %T SET %T %T"
	}

	cmd, err := db.Exec(ctx,
		query,
		table,
		vals,
		onConflict,
	)

	if err == nil {
		vals.reset()
		count, err = cmd.RowsAffected()
	}

	return
}

func DoUpdate(ignoreColumns ...string) func(vals *Values) EncodeQuery {
	return func(vals *Values) EncodeQuery {
		return func(buf *buffer.Buffer, queryArgs *[]any) {
			var written bool

			for i := range vals.columns {
				if slices.Contains(ignoreColumns, vals.columns[i]) {
					continue
				}

				if written {
					buf.WriteString(", ")
				} else {
					buf.WriteString("ON DUPLICATE KEY UPDATE ")
					written = true
				}

				writeIdentifier(buf, vals.columns[i])
				buf.WriteString(" = VALUES(")
				writeIdentifier(buf, vals.columns[i])
				buf.WriteString(")")
			}
		}
	}
}
