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
	sizerSrvc *SizerService
}

// NewPacketsService is a constructor of the PacketsService.
func NewPacketsService(sizerSrvc *SizerService) *PacketsService {
	return &PacketsService{
		sizerSrvc: sizerSrvc,
	}
}

// GetPackets ...
func (packets PacketsService) GetPackets(ctx context.Context, itemsToPack int) (map[int]int, error) {
	if itemsToPack <= 0 {
		slog.ErrorContext(ctx,
			ErrorNegativeOrZeroItems,
			"incoming_items", itemsToPack)
		return map[int]int{}, errors.New(ErrorNegativeOrZeroItems)
	}

	return getMinNecessaryPacks(itemsToPack, packets.sizerSrvc.SortedSizes), nil
}

// getMinNecessaryPacks calculates minimum packs quantity for given items based on packs sizes.
func getMinNecessaryPacks(items int, sortedSizes []int) map[int]int {
	necessaryPacks := make(map[int]int)
	lastUsedPackIndex := len(sortedSizes) - 1

	diff := 0
	for lastUsedPackIndex > 0 {
		if items-sortedSizes[lastUsedPackIndex] >= 0 {
			necessaryPacks[sortedSizes[lastUsedPackIndex]]++
			items -= sortedSizes[lastUsedPackIndex]
		} else {
			if _, exists := necessaryPacks[sortedSizes[lastUsedPackIndex]]; exists {
				diff = sortedSizes[lastUsedPackIndex] - items
				if sortedSizes[lastUsedPackIndex-1] > diff {
					necessaryPacks[sortedSizes[lastUsedPackIndex]]++
					items -= sortedSizes[lastUsedPackIndex]
					break
				}
			}
			lastUsedPackIndex--
		}
	}

	if items > 0 {
		for _, size := range sortedSizes {
			if size >= items {
				necessaryPacks[size]++
				break
			}
		}
	}

	return necessaryPacks
}
