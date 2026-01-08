package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func Json(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, r *http.Request, status int, data any) {
	pb := ProblemDetails{
		Instance: r.URL.String(),
		Type:     "about:blank",
		Title:    http.StatusText(status),
		Detail:   fmt.Sprintf("%s", data),
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(pb)
}

func Header(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
