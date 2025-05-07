package mysql

import "github.com/webmafia/fast/buffer"

var _ QueryEncoder = Cond(nil)

type Cond func(buf *buffer.Buffer, queryArgs *[]any)

func (c Cond) EncodeQuery(buf *buffer.Buffer, queryArgs *[]any) {
	c(buf, queryArgs)
}
