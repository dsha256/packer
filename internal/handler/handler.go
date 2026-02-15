package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/dsha256/packer/internal/middleware"
	"github.com/dsha256/packer/internal/packer"
	"github.com/dsha256/packer/internal/responder"
	"github.com/dsha256/packer/pkg/cache"
)

var ErrMethodNotAllowed = errors.New("method not allowed")

type Handler struct {
	logger *slog.Logger
	packer packer.Packer
	cache  cache.Cache
}

func New(
	logger *slog.Logger,
	packer packer.Packer,
	cache cache.ClosableCache,
) *Handler {
	return &Handler{
		logger: logger,
		packer: packer,
		cache:  cache,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/packet/calculate", h.wrapHandler(h.handlePacketsCalculation))
	mux.Handle("/api/v1/packet/size", h.wrapHandler(h.handlePacketSizes))
	mux.Handle("/api/v1/health", h.wrapHandler(h.handleHealth))
	h.logger.Info("Routes registered")
}

func (h *Handler) wrapHandler(handler http.HandlerFunc) http.Handler {
	return middleware.LoggingMiddleware(
		h.logger,
		middleware.RecoverMiddleware(
			h.logger,
			handler,
		),
	)
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	responder.WriteSuccess(w, http.StatusOK, "All services are up and running", json.RawMessage{})
}

func (h *Handler) handleError(w http.ResponseWriter, err error, status int) {
	h.logger.Error("Error handling request", "error", err)
	responder.WriteError(w, status, err)
}
