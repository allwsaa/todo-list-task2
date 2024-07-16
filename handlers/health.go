package handlers

import (
	"net/http"
)

// HealthCheckHandler handles the health check request.
// @Summary Health check
// @Description Check if the service is running.
// @Success 200 {string} string "OK"
// @Router /health [get]
func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
