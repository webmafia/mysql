package mysql

import (
	"github.com/webmafia/fast/buffer"
)

var _ StringEncoder = ChainedIdentifier{Identifier(""), Identifier("")}

type ChainedIdentifier [2]StringEncoder

// EncodeString implements StringEncoder.
func (t ChainedIdentifier) EncodeString(b *buffer.Buffer) {
	t[0].EncodeString(b)
	b.WriteByte('.')
	t[1].EncodeString(b)
}

func (t ChainedIdentifier) Col(col string) ChainedIdentifier {
	return ChainedIdentifier{t, Identifier(col)}
}

func (t ChainedIdentifier) Alias(col string) Alias {
	return Alias{t[1], Identifier(col)}
}
