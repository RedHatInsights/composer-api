package workspace

type Handler struct {
	// Future: DB, cache, etc.
}

func New() *Handler {
	return &Handler{}
}
