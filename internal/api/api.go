package api

import (
	"abb-exporter/internal/api/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func BindApi(router *chi.Mux, ah *handler.Meter) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api", func(r chi.Router) {
		bindAbbApi(r, ah)
		bindHealthApi(r)
	})
}

func bindHealthApi(router chi.Router) {
	sh := handler.NewStatus()
	router.Get("/health", sh.Health)
	router.Get("/status", sh.Status)
}

func bindAbbApi(router chi.Router, h *handler.Meter) {
	router.Route("/meter/{name}", func(r chi.Router) {
		r.Get("/info", h.ReadInfo)
		r.Get("/usage-status", h.ReadUsageStatus)
		r.Get("/total-active-power", h.ReadTotalActivePower)
	})
}
