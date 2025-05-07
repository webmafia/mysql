package mysql

import (
	"strings"

	"github.com/webmafia/fast"
	"github.com/webmafia/fast/buffer"
)

func writeAny(b *buffer.Buffer, args *[]any, val any) {
	switch v := fast.Noescape(val).(type) {

	case StringEncoder:
		if v != nil {
			v.EncodeString(b)
		}

	case QueryEncoder:
		if v != nil {
			v.EncodeQuery(b, args)
		}

	default:
		writeQueryArg(b, args, val)

	}
}

func writeQueryArg(b *buffer.Buffer, args *[]any, val any) {
	*args = append(*fast.Noescape(args), fast.Noescape(val))
	b.WriteByte('$')
	b.Str().WriteInt(len(*args))
}

func writeAnyIdentifier(b *buffer.Buffer, str any) {
	switch v := fast.Noescape(str).(type) {
	case StringEncoder:
		v.EncodeString(b)
	case string:
		dot := strings.IndexByte(v, '.')

		if dot >= 0 {
			writeIdentifier(b, v[:dot])
			b.WriteByte('.')
			v = v[dot+1:]
		}

		writeIdentifier(b, v)
	}
}

func writeIdentifier(b *buffer.Buffer, str string) {
	b.WriteByte('"')
	b.WriteString(str)
	b.WriteByte('"')
}
