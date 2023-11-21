package server

import (
	"github.com/dsha256/packer/internal/validator"
)

func (s *Server) validateSizeOnValue(v *validator.Validator, size int) {
	v.Check(size > 0, "size", "size must be positive number")
}
