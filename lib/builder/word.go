package builder

type word struct {
	w string
}

func Word(w string) Builder {
	if w == "" {
		return &null{}
	}
	return &word{w}
}

func (b *word) Build() string {
	return b.w
}
