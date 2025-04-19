package handler

import (
	"net/http"

	"github.com/goccy/go-json"

	"github.com/dsha256/packer/internal/responder"
	"github.com/dsha256/packer/internal/types"
	"github.com/dsha256/packer/internal/validation"
)

func (h *Handler) handlePacketSizes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListPacketSizes(w, r)
	case http.MethodPut:
		h.handlePutPacketSizes(w, r)
	default:
		h.handleError(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListPacketSizes(w http.ResponseWriter, r *http.Request) {
	packetSizes, err := h.packer.ListPacketSizes(r.Context())
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)

		return
	}

	responder.WriteSuccess(w, http.StatusOK, "", map[string][]types.PacketSize{
		"packet_sizes": packetSizes,
	})
}

type PutPacketSizesRequest struct {
	Sizes []types.PacketSize `json:"sizes"`
}

func (h *Handler) handlePutPacketSizes(w http.ResponseWriter, r *http.Request) {
	var sizes PutPacketSizesRequest
	if err := json.NewDecoder(r.Body).Decode(&sizes); err != nil {
		h.handleError(w, err, http.StatusBadRequest)

		return
	}

	if err := validation.ValidatePacketSizes(sizes.Sizes); err != nil {
		h.handleError(w, err, http.StatusBadRequest)

		return
	}

	if err := h.packer.SetPacketSizes(r.Context(), sizes.Sizes); err != nil {
		return
	}

	responder.WriteSuccess(w, http.StatusOK, "Packet sizes have been put successfully", json.RawMessage{})
}
