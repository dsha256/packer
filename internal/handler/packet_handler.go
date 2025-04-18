package handler

import (
	"net/http"

	"github.com/goccy/go-json"

	"github.com/dsha256/packer/internal/responder"
)

func (h *Handler) handlePackets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetPacket(w, r)
	default:
		h.handleError(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGetPacket(w http.ResponseWriter, r *http.Request) {
	responder.WriteSuccess(w, http.StatusOK, "handleGetPacket is up and running", json.RawMessage{})
}
