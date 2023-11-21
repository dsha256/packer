package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (s *Server) logError(r *http.Request, err error) {
	slog.Error(err.Error(),
		slog.Any("request_method", r.Method),
		slog.Any("request_url", r.URL.String()),
	)
}

func (s *Server) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}
	err := s.writeJSON(w, status, env, nil)
	if err != nil {
		s.logError(r, err)
		w.WriteHeader(500)
	}
}

func (s *Server) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	s.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (s *Server) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	s.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (s *Server) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (s *Server) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	s.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (s *Server) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	s.errorResponse(w, r, http.StatusNotFound, message)
}

func (s *Server) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	s.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
