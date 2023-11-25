package server

import (
	"net/http"

	"github.com/dsha256/packer/internal/validator"
)

func (s *Server) listSizesHandler(w http.ResponseWriter, r *http.Request) {
	err := s.writeJSON(w, http.StatusOK, envelope{"sizes": s.SizerSrvc.ListSizes()}, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}

func (s *Server) addSizeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Size int `json:"size"`
	}

	err := s.readJSON(w, r, &input)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	s.validateSizeOnValue(v, input.Size)
	if !v.Valid() {
		s.failedValidationResponse(w, r, v.Errors)
		return
	}

	sizes, err := s.SizerSrvc.AddSize(r.Context(), input.Size)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	err = s.writeJSON(w, http.StatusOK, envelope{
		"sorted_sizes": sizes,
	}, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}

func (s *Server) putSizesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Sizes []int `json:"sizes"`
	}

	err := s.readJSON(w, r, &input)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	for _, size := range input.Sizes {
		s.validateSizeOnValue(v, size)
	}
	if !v.Valid() {
		s.failedValidationResponse(w, r, v.Errors)
		return
	}

	sizes, err := s.SizerSrvc.PutSizes(r.Context(), input.Sizes)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	err = s.writeJSON(w, http.StatusOK, envelope{
		"sorted_sizes": sizes,
	}, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}

func (s *Server) deleteSizeHandler(w http.ResponseWriter, r *http.Request) {
	size, err := s.readSizeParam(r)
	if err != nil {
		s.notFoundResponse(w, r)
		return
	}

	v := validator.New()
	s.validateSizeOnValue(v, size)
	if !v.Valid() {
		s.failedValidationResponse(w, r, v.Errors)
		return
	}

	sizes, err := s.SizerSrvc.DeleteSize(r.Context(), size)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	err = s.writeJSON(w, http.StatusOK, envelope{
		"sorted_sizes": sizes,
	}, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}
