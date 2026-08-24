package probe

type Handler struct {
	// Future: DB pool for health checks, etc.
}

func New() *Handler {
	return &Handler{}
}
