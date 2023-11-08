package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", s.healthcheckHandler)

	return s.metrics(s.recoverPanic(s.enableCORS(s.rateLimit(router))))
}
