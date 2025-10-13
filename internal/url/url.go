package url

// URL represents a URL with an ID and original string.
type URL struct {
	ID       string
	Original string
}

// NewURL creates a new URL with the given ID and original string.
func NewURL(id, original string) *URL {
	return &URL{ID: id, Original: original}
}
