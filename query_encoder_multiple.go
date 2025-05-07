package mysql

import (
	"github.com/webmafia/fast/buffer"
)

func Multiple(enc ...QueryEncoder) EncodeQuery {
	return func(buf *buffer.Buffer, queryArgs *[]any) {
		for i := range enc {
			if i != 0 {
				buf.WriteString(", ")
			}

			enc[i].EncodeQuery(buf, queryArgs)
		}
	}
}
