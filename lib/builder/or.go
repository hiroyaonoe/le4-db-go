package builder

import "strings"

type or struct {
	b []Builder
}

func Or(b ...Builder) Builder {
	for i, v := range b {
		if _, ok := v.(*null); ok {
			return Or(append(b[:i], b[i+1:]...)...)
		}
	}
	if len(b) == 0 {
		return &null{}
	}
	if len(b) == 1 {
		return b[0]
	}
	return &or{b}
}

func (b *or) Build() string {
	s := make([]string, len(b.b))
	for i, v := range b.b {
		s[i] = v.Build()
	}
	return "( " + strings.Join(s, " OR ") + " )"
}
