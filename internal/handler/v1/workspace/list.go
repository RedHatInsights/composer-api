package workspace

import (
	"net/http"

	"github.com/RedHatInsights/composer-api/internal/response"
)

type ListAPIResponse struct {
	Workspaces []any `json:"workspaces"`
}

// TODO: stub — returns an empty list; implementation pending.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, ListAPIResponse{Workspaces: []any{}})
}
