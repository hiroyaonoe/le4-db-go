package builder

type where struct {
	s Builder
	q Builder
}

func Where(s Builder, q Builder) Builder {
	if _, ok := q.(*null); ok {
		return s
	}
	return &where{
		s: s,
		q: q,
	}
}

func (b where) Build() string {
	return b.s.Build() + " WHERE " + b.q.Build()
}
