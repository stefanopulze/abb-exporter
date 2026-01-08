package handler

import (
	"abb-exporter/internal/api/response"
	"abb-exporter/internal/meter"
	"errors"
	"fmt"
	"net/http"
)

func NewMeter(client *meter.Group) *Meter {
	return &Meter{
		client: client,
	}
}

type Meter struct {
	client *meter.Group
}

func (h Meter) ReadTotalActivePower(w http.ResponseWriter, r *http.Request) {
	client, err := h.meterFromPathName(r)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err)
		return
	}

	activePower, err := client.QueryTotalActiveImport()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Json(w, http.StatusOK, map[string]any{"activePower": activePower})
}

func (h Meter) ReadUsageStatus(w http.ResponseWriter, r *http.Request) {
	client, err := h.meterFromPathName(r)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err)
		return
	}
	regs, err := client.QueryUsageStatus()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Json(w, http.StatusOK, regs)
}

func (h Meter) ReadInfo(w http.ResponseWriter, r *http.Request) {
	client, err := h.meterFromPathName(r)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err)
		return
	}

	info, err := client.QueryInfo()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Json(w, http.StatusOK, info)
}

func (h Meter) meterFromPathName(r *http.Request) (meter.Meter, error) {
	name := r.PathValue("name")
	if len(name) < 3 {
		return nil, errors.New("invalid meter name")
	}

	m, found := h.client.Get(name)
	if !found {
		return nil, fmt.Errorf("meter with name %s not found", name)
	}

	return m, nil
}
