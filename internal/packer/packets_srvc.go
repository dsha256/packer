package packer

import (
	"context"
	"errors"
	"log/slog"
)

const (
	// ErrorNegativeOrZeroItems ...
	ErrorNegativeOrZeroItems = "items must be more than 0"
)

// Ensure PacketsService defined types fully satisfy Packer interfaces.
var _ Packer = &PacketsService{}

// PacketsService holds the Packets service related params.
type PacketsService struct {
}

// NewPacketsService is a constructor of the PacketsService.
func NewPacketsService() *PacketsService {
	return &PacketsService{}
}

// GetPackets ...
func (packets PacketsService) GetPackets(ctx context.Context, itemsToPack int) (map[int]int, error) {
	if itemsToPack <= 0 {
		slog.ErrorContext(ctx,
			ErrorNegativeOrZeroItems,
			"incoming_items", itemsToPack)
		return map[int]int{}, errors.New(ErrorNegativeOrZeroItems)
	}

	return getMinNecessaryPacks(itemsToPack), nil
}

// getMinNecessaryPacks calculates minimum packs quantity for given items based on packs sizes.
func getMinNecessaryPacks(items int) map[int]int {
	necessaryPacks := make(map[int]int)
	lastUsedPackIndex := len(SortedSizes) - 1

	for lastUsedPackIndex > 0 {
		if items-SortedSizes[lastUsedPackIndex] >= 0 {
			necessaryPacks[SortedSizes[lastUsedPackIndex]]++
			items -= SortedSizes[lastUsedPackIndex]
		} else {
			lastUsedPackIndex--
		}
	}

	if items > 0 {
		for _, size := range SortedSizes {
			if size >= items {
				necessaryPacks[size]++
				break
			}
		}
	}

	return necessaryPacks
}
