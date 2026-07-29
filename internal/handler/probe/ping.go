package probe

import (
	"net/http"

	"github.com/RedHatInsights/composer-api/internal/response"
)

type PingAPIResponse struct {
	Ping string `json:"ping"`
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, PingAPIResponse{Ping: "pong"})
}
