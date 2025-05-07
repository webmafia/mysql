package mysql

import (
	"github.com/webmafia/fast"
	"github.com/webmafia/fast/buffer"
)

func Raw(s string, args ...any) QueryEncoder {
	return Cond(func(buf *buffer.Buffer, queryArgs *[]any) {
		if len(args) > 0 {
			buf.Str().WritefCb(s, fast.Noescape(args), func(b *buffer.Buffer, c byte, v any) error {
				writeAny(b, queryArgs, v)
				return nil
			})
		} else {
			buf.WriteString(s)
		}

	})
}
