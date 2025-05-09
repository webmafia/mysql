package mysql

import (
	"github.com/webmafia/fast/buffer"
)

var _ StringEncoder = Identifier("")

type Identifier string

// EncodeString implements StringEncoder.
func (t Identifier) EncodeString(b *buffer.Buffer) {
	writeIdentifier(b, string(t))
}

func (t Identifier) Col(col string) ChainedIdentifier {
	return ChainedIdentifier{t, Identifier(col)}
}

func (t Identifier) Alias(col string) Alias {
	return Alias{t, Identifier(col)}
}

func writeIdentifiers(b *buffer.Buffer, strs []string) {
	for i := range strs {
		if i != 0 {
			b.WriteByte(',')
		}

		writeIdentifier(b, strs[i])
	}
}
