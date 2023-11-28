package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(s.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(s.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", s.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/api/v1/sizes", s.listSizesHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/sizes", s.addSizeHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/sizes", s.putSizesHandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/sizes/:size", s.deleteSizeHandler)

	router.HandlerFunc(http.MethodGet, "/api/v1/docs", s.docsHandler)

	return s.metrics(s.recoverPanic(s.enableCORS(s.rateLimit(router))))
}
