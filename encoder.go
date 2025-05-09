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

var _ StringEncoder = (EncodeString)(nil)

type EncodeString func(buf *buffer.Buffer)

func (fn EncodeString) EncodeString(buf *buffer.Buffer) {
	fn(buf)
}

var _ QueryEncoder = (EncodeQuery)(nil)

type EncodeQuery func(buf *buffer.Buffer, queryArgs *[]any)

func (fn EncodeQuery) EncodeQuery(buf *buffer.Buffer, queryArgs *[]any) {
	fn(buf, queryArgs)
}

var queryEncoderNoop = EncodeQuery(func(_ *buffer.Buffer, _ *[]any) {})
