package mysql

import (
	"github.com/webmafia/fast/buffer"
)

var _ StringEncoder = Alias{Identifier(""), Identifier("")}

type Alias [2]StringEncoder

// EncodeString implements StringEncoder.
func (t Alias) EncodeString(b *buffer.Buffer) {
	t[0].EncodeString(b)
	b.WriteString(" AS ")
	t[1].EncodeString(b)
}

func (t Alias) Col(col string) ChainedIdentifier {
	return ChainedIdentifier{t[1], Identifier(col)}
}
