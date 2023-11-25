package server

import (
	"net/http"

	"github.com/dsha256/packer/internal/validator"
)

func (s *Server) getPacksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Items int `json:"items"`
	}

	err := s.readJSON(w, r, &input)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	s.validateItemsOnValue(v, input.Items)
	if !v.Valid() {
		s.failedValidationResponse(w, r, v.Errors)
		return
	}

	packets, err := s.PackerSrvc.GetPackets(r.Context(), input.Items)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	err = s.writeJSON(w, http.StatusOK, envelope{
		"packets": packets,
	}, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}
