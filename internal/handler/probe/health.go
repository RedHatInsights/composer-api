package probe

import (
	"net/http"

	"github.com/RedHatInsights/composer-api/internal/response"
)

type HealthAPIResponse struct {
	Healthy bool `json:"healthy"`
}

// TODO: expand to check dependencies (database, external services, etc.)
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, HealthAPIResponse{Healthy: true})
}
