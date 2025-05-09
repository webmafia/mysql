package mysql

import "github.com/webmafia/fast/buffer"

var (
	_ QueryEncoder = (*Values)(nil)
	_ QueryEncoder = cols{}
	_ QueryEncoder = vals{}
)

//go:inline
func (db *DB) AcquireValues() *Values {
	vals, ok := db.valPool.Get().(*Values)

	if ok {
		return vals
	}

	return &Values{
		columns: make([]string, 0, 8),
		values:  make([]any, 0, 8),
	}
}

//go:inline
func (db *DB) ReleaseValues(vals *Values) {
	vals.reset()
	db.valPool.Put(vals)
}

type Values struct {
	columns []string
	values  []any
}

func (r *Values) reset() {
	clear(r.columns)
	clear(r.values)
	r.columns = r.columns[:0]
	r.values = r.values[:0]
}

func (r *Values) Value(column string, value any) *Values {
	r.columns = append(r.columns, column)
	r.values = append(r.values, value)

	return r
}

func (r *Values) Len() int {
	return len(r.columns)
}

func (r *Values) Empty() bool {
	return len(r.columns) == 0
}

// EncodeQuery implements QueryEncoder.
func (r *Values) EncodeQuery(buf *buffer.Buffer, queryArgs *[]any) {
	for i := range r.columns {
		if i != 0 {
			buf.WriteString(", ")
		}

		writeIdentifier(buf, r.columns[i])
		buf.WriteString(" = ")
		writeQueryArg(buf, queryArgs, r.values[i])
	}
}

func (r *Values) colEncoder() QueryEncoder {
	return cols{v: r}
}

type cols struct {
	v *Values
}

// EncodeQuery implements QueryEncoder.
func (c cols) EncodeQuery(buf *buffer.Buffer, _ *[]any) {
	writeIdentifiers(buf, c.v.columns)
}

func (r *Values) valEncoder() QueryEncoder {
	return vals{v: r}
}

type vals struct {
	v *Values
}

// EncodeQuery implements QueryEncoder.
func (c vals) EncodeQuery(buf *buffer.Buffer, args *[]any) {
	for i := range c.v.values {
		if i != 0 {
			buf.WriteByte(',')
		}

		writeQueryArg(buf, args, c.v.values[i])
	}
}
