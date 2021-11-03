package builder

import "strings"

type and struct {
	b []Builder
}

func And(b ...Builder) Builder {
	for i, v := range b {
		if _, ok := v.(*null); ok {
			return And(append(b[:i], b[i+1:]...)...)
		}
	}
	if len(b) == 0 {
		return &null{}
	}
	if len(b) == 1 {
		return b[0]
	}
	return &and{b}
}

func (b *and) Build() string {
	s := make([]string, len(b.b))
	for i, v := range b.b {
		s[i] = v.Build()
	}
	return "( " + strings.Join(s, " AND ") + " )"
}
