package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/dsha256/packer/internal/responder"
	"github.com/dsha256/packer/pkg/cache"
	"github.com/dsha256/packer/pkg/safeconv"
)

var (
	ErrInvalidItems  = errors.New("items should be positive integer")
	ErrItemsTooLarge = errors.New("items exceed maximum allowed value of 1 billion")
)

const MaxAllowedItems = 1_000_000_000

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
		h.logger.Warn("Error parsing the request payload", "err", err)
		h.handleError(w, err, http.StatusBadRequest)

		return
	}

	items := r.Form.Get("items")
	itemsInt := safeconv.ParseInt(items)
	if itemsInt < 1 {
		h.logger.Warn("Invalid incoming items", "err", ErrInvalidItems)
		h.handleError(w, ErrInvalidItems, http.StatusBadRequest)

		return
	}

	if itemsInt > MaxAllowedItems {
		h.logger.Warn("Items exceed maximum allowed", "items", itemsInt, "max", MaxAllowedItems)
		h.handleError(w, ErrItemsTooLarge, http.StatusBadRequest)

		return
	}

	cachedPackets, err := h.cache.Get(r.Context(), items)
	if err != nil {
		if errors.Is(err, cache.ErrNoKey) {
			h.logger.Info("Items is not cached", "err", err)
		} else {
			h.logger.Error("Failed to get items", "err", err)
			h.handleError(w, err, http.StatusInternalServerError)

			return
		}
	}
	if err == nil {
		responder.WriteSuccess(w, http.StatusOK, "", map[string]any{
			"optimal_packets": cachedPackets,
		})

		return
	}

	packets, err := h.packer.GetOptimalPackets(r.Context(), itemsInt)
	if err != nil {
		h.logger.Error("Failed to get optimal packets", "err", err)
		h.handleError(w, err, http.StatusInternalServerError)

		return
	}

	if err = h.cache.Set(r.Context(), items, packets, 1*time.Hour); err != nil {
		h.logger.Error("Failed to set items to cache", "err", err)
		h.handleError(w, err, http.StatusInternalServerError)

		return
	}

	responder.WriteSuccess(w, http.StatusOK, "", map[string]any{
		"optimal_packets": packets,
	})
}
