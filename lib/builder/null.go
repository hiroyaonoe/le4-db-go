package builder

type null struct {
}

func Null() Builder {
	return &null{}
}
func (b *null) Build() string {
	return ""
}
