package mysql

import (
	"github.com/webmafia/fast/buffer"
)

var _ StringEncoder = Identifier("")

type Identifier string

// EncodeString implements StringEncoder.
func (t Identifier) EncodeString(b *buffer.Buffer) {
	b.WriteByte('`')
	b.WriteString(string(t))
	b.WriteByte('`')
}

func (t Identifier) Col(col string) ChainedIdentifier {
	return ChainedIdentifier{t, Identifier(col)}
}

func (t Identifier) Alias(col string) Alias {
	return Alias{t, Identifier(col)}
}
