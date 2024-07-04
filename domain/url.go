package domain

type URL struct {
	ID       string
	Original string
}

func NewURL(id, original string) *URL {
	return &URL{ID: id, Original: original}
}
