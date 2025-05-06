package mysql

import (
	"github.com/webmafia/fast/buffer"
)

type StringEncoder interface {
	EncodeString(buf *buffer.Buffer)
}

type QueryEncoder interface {
	EncodeQuery(buf *buffer.Buffer, queryArgs *[]any)
}
