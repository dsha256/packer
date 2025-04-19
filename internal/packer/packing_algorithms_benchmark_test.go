package packer_test

import (
	"testing"

	"github.com/dsha256/packer/internal/packer"
	"github.com/dsha256/packer/internal/types"
)

func benchmarkCalculateOptimalPacketsForItems(b *testing.B, calculateFunc func(params *packer.CalculateOptimalPacketsForItemsParams) map[types.PacketSize]types.PacketQuantity) {
	testCases := []struct {
		Name        string
		PacketSizes []types.PacketSize
		Items       int
	}{
		{
			Name:        "ExtraSmallQuantityOfItems",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       250,
		},
		{
			Name:        "SmallQuantityOfItems",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       12001,
		},
		{
			Name:        "MediumQuantityOfItems",
			PacketSizes: []types.PacketSize{23, 31, 53},
			Items:       100_000,
		},
		{
			Name:        "LargeQuantityOfItems",
			PacketSizes: []types.PacketSize{23, 31, 53},
			Items:       500_000,
		},
		{
			Name:        "ExtraLargeQuantityOfItems",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       10_000_000,
		},
	}

	for _, testCase := range testCases {
		b.Run(testCase.Name, func(b *testing.B) {
			params := &packer.CalculateOptimalPacketsForItemsParams{
				Items:       testCase.Items,
				PacketSizes: testCase.PacketSizes,
			}

			b.ResetTimer()

			for b.Loop() {
				calculateFunc(params)
			}
		})
	}
}

func Benchmark_CalculateOptimalPacketsForItemsV2(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItems(b, packer.CalculateOptimalPacketsForItemsV2)
}

func Benchmark_CalculateOptimalPacketsForItemsV1(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItems(b, packer.CalculateOptimalPacketsForItemsV1)
}
