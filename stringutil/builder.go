package stringutil

import (
	"strings"

	"github.com/spf13/cast"
)

// Builder 支持连写
type Builder struct {
	*strings.Builder
}

func NewBuilder() *Builder {
	return &Builder{Builder: &strings.Builder{}}
}

func (b *Builder) Write(a any) *Builder {
	if b.Builder == nil {
		b.Builder = &strings.Builder{}
	}

	switch v := a.(type) {
	case []byte:
		b.Builder.Write(v)
	case byte:
		b.Builder.WriteByte(v)
	case rune:
		b.Builder.WriteRune(v)
	case string:
		b.Builder.WriteString(v)
	default:
		b.Builder.WriteString(cast.ToString(v))
	}

	return b
}
