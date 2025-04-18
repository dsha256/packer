package handler

import (
	"errors"
	"net/http"

	"github.com/dsha256/packer/internal/responder"
	"github.com/dsha256/packer/pkg/safeconv"
)

var (
	ErrInvalidItems = errors.New("items should be positive integer")
)

func (h *Handler) handlePacketsCalculation(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetOptimalPackets(w, r)
	default:
		h.handleError(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGetOptimalPackets(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	items := r.Form.Get("items")
	itemsInt := safeconv.ParseInt(items)
	if itemsInt < 1 {
		h.handleError(w, ErrInvalidItems, http.StatusBadRequest)
		return
	}

	packets, err := h.packer.GetOptimalPackets(r.Context(), itemsInt)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "", map[string]any{
		"optimal_packets": packets,
	})
}
