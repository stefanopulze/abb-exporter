package handler

import (
	"abb-exporter/internal/api/response"
	"net/http"
)

func NewStatus() *StatusHandler {
	return &StatusHandler{}
}

type StatusHandler struct {
}

func (sh StatusHandler) Health(w http.ResponseWriter, _ *http.Request) {
	response.Header(w, http.StatusOK)
}

func (sh StatusHandler) Status(w http.ResponseWriter, _ *http.Request) {
	data := map[string]any{
		"status": "ok",
	}

	response.Json(w, http.StatusOK, data)
}
