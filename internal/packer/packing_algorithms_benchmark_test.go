package packer

import (
	"testing"

	"github.com/dsha256/packer/internal/types"
)

func benchmarkCalculateOptimalPacketsForItems(b *testing.B, calculateFunc func(params *CalculateOptimalPacketsForItemsParams) map[types.PacketSize]types.PacketQuantity) {
	testCases := []struct {
		Name        string
		Items       int
		PacketSizes []types.PacketSize
	}{
		{"ExtraSmallQuantityOfItems", 250, []types.PacketSize{250, 500, 1000, 2000, 5000}},
		{"SmallQuantityOfItems", 12001, []types.PacketSize{250, 500, 1000, 2000, 5000}},
		{"MediumQuantityOfItems", 100_000, []types.PacketSize{23, 31, 53}},
		{"LargeQuantityOfItems", 500_000, []types.PacketSize{23, 31, 53}},
		{"ExtraLargeQuantityOfItems", 10_000_000, []types.PacketSize{250, 500, 1000, 2000, 5000}},
	}

	for _, testCase := range testCases {
		b.Run(testCase.Name, func(b *testing.B) {
			params := &CalculateOptimalPacketsForItemsParams{
				Items:       testCase.Items,
				PacketSizes: testCase.PacketSizes,
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				calculateFunc(params)
			}
		})
	}
}

func Benchmark_CalculateOptimalPacketsForItemsV2(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItems(b, CalculateOptimalPacketsForItemsV2)
}

func Benchmark_CalculateOptimalPacketsForItemsV1(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItems(b, CalculateOptimalPacketsForItemsV1)
}
