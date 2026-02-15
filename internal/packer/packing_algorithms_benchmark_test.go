package packer_test

import (
	"testing"

	"github.com/dsha256/packer/internal/packer"
	"github.com/dsha256/packer/internal/types"
)

//nolint:dupl // Clearer in this case.
func benchmarkCalculateOptimalPacketsForItemsWithProductOfTenSizes(b *testing.B, calculateFunc func(params *packer.CalculateOptimalPacketsForItemsParams) map[types.PacketSize]types.PacketQuantity) {
	b.Helper()

	testCases := []struct {
		Name        string
		PacketSizes []types.PacketSize
		Items       int
	}{
		{
			Name:        "~50k_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       50_123,
		},
		{
			Name:        "~100k_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       100_123,
		},
		{
			Name:        "~250k_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       250_123,
		},
		{
			Name:        "~500k_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       500_123,
		},
		{
			Name:        "~1M_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       1_123_123,
		},
		{
			Name:        "~10M_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       10_123_123,
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

//nolint:dupl // Clearer in this case.
func benchmarkCalculateOptimalPacketsForItemsWithPrimeSizes(b *testing.B, calculateFunc func(params *packer.CalculateOptimalPacketsForItemsParams) map[types.PacketSize]types.PacketQuantity) {
	b.Helper()

	testCases := []struct {
		Name        string
		PacketSizes []types.PacketSize
		Items       int
	}{
		{
			Name:        "~50k_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       50_123,
		},
		{
			Name:        "~100k_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       100_123,
		},
		{
			Name:        "~250k_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       250_123,
		},
		{
			Name:        "~500k_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       500_123,
		},
		{
			Name:        "~1M_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       1_123_123,
		},
		{
			Name:        "~10M_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       10_123_123,
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

func Benchmark_CalculateOptimalPacketsForItemsV2_ProductOfTenSizes(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItemsWithProductOfTenSizes(b, packer.CalculateOptimalPacketsForItemsV2)
}

func Benchmark_CalculateOptimalPacketsForItemsV1_ProductOfTenSizes(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItemsWithProductOfTenSizes(b, packer.CalculateOptimalPacketsForItemsV1)
}

func Benchmark_CalculateOptimalPacketsForItemsV2_PrimeSizes(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItemsWithPrimeSizes(b, packer.CalculateOptimalPacketsForItemsV2)
}

func Benchmark_CalculateOptimalPacketsForItemsV1_PrimeSizes(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItemsWithPrimeSizes(b, packer.CalculateOptimalPacketsForItemsV1)
}
