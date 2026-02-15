package validation

import (
	"errors"

	"github.com/dsha256/packer/internal/types"
)

var (
	ErrNonPositiveSize = errors.New("size should be a positive integer")
	ErrDuplicatedSizes = errors.New("sizes should be unique")
)

func ValidatePacketSizes(sizes []types.PacketSize) error {
	tempSizes := make(map[types.PacketSize]types.PacketSize, len(sizes))
	for _, size := range sizes {
		if size < 1 {
			return ErrNonPositiveSize
		}
		tempSizes[size]++
	}

	for _, frequency := range tempSizes {
		if frequency > 1 {
			return ErrDuplicatedSizes
		}
	}

	return nil
}
