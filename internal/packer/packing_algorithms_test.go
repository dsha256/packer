package packer_test

import (
	"reflect"
	"testing"

	"github.com/dsha256/packer/internal/packer"
	"github.com/dsha256/packer/internal/types"
)

func Test_CalculateOptimalPacketsForItems(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedOptimalPacks map[types.PacketSize]types.PacketQuantity
		Name                 string
		PacketSizes          []types.PacketSize
		Items                int
	}{
		{
			Name:        "1",
			Items:       1,
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			ExpectedOptimalPacks: map[types.PacketSize]types.PacketQuantity{
				250: 1,
			},
		},
		{
			Name:        "250",
			Items:       250,
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			ExpectedOptimalPacks: map[types.PacketSize]types.PacketQuantity{
				250: 1,
			},
		},
		{
			Name:        "251",
			Items:       251,
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			ExpectedOptimalPacks: map[types.PacketSize]types.PacketQuantity{
				500: 1,
			},
		},
		{
			Name:        "501",
			Items:       501,
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			ExpectedOptimalPacks: map[types.PacketSize]types.PacketQuantity{
				500: 1,
				250: 1,
			},
		},
		{
			Name:        "12001",
			Items:       12001,
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			ExpectedOptimalPacks: map[types.PacketSize]types.PacketQuantity{
				5000: 2,
				2000: 1,
				250:  1,
			},
		},
		{
			Name:        "500000",
			Items:       500000,
			PacketSizes: []types.PacketSize{23, 31, 53},
			ExpectedOptimalPacks: map[types.PacketSize]types.PacketQuantity{
				23: 2,
				31: 7,
				53: 9429,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			result := packer.CalculateOptimalPacketsForItemsV1(&packer.CalculateOptimalPacketsForItemsParams{
				Items:       testCase.Items,
				PacketSizes: testCase.PacketSizes,
			})
			if !reflect.DeepEqual(result, testCase.ExpectedOptimalPacks) {
				t.Errorf("CalculateOptimalPacketsForItemsV1: Expected: %v \nGot: %v", testCase.ExpectedOptimalPacks, result)
			}
		})
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			result := packer.CalculateOptimalPacketsForItemsV2(&packer.CalculateOptimalPacketsForItemsParams{
				Items:       testCase.Items,
				PacketSizes: testCase.PacketSizes,
			})
			if !reflect.DeepEqual(result, testCase.ExpectedOptimalPacks) {
				t.Errorf("CalculateOptimalPacketsForItemsV2: Expected: %v \nGot: %v", testCase.ExpectedOptimalPacks, result)
			}
		})
	}
}
