package packer

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func newPacker() *PacketsService {
	return NewPacketsService()
}

func TestPacketsService_GetPackets(t *testing.T) {
	packer := newPacker()

	packets, err := packer.GetPackets(context.Background(), 0)
	require.True(t, reflect.DeepEqual(packets, map[int]int{}))
	require.Error(t, err)
	require.Equal(t, err.Error(), ErrorNegativeOrZeroItems)

	responseFor10Items := map[int]int{250: 1}
	packets, err = packer.GetPackets(context.Background(), 10)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(packets, responseFor10Items))
}

func Test_getMinNecessaryPacks(t *testing.T) {
	testCases := []struct {
		Items              int
		WantNecessaryPacks map[int]int
	}{
		{
			Items: 1,
			WantNecessaryPacks: map[int]int{
				250: 1,
			},
		},
		{
			Items: 250,
			WantNecessaryPacks: map[int]int{
				250: 1,
			},
		},
		{
			Items: 251,
			WantNecessaryPacks: map[int]int{
				500: 1,
			},
		},
		{
			Items: 501,
			WantNecessaryPacks: map[int]int{
				500: 1,
				250: 1,
			},
		},
		{
			Items: 12001,
			WantNecessaryPacks: map[int]int{
				5000: 2,
				2000: 1,
				250:  1,
			},
		},
	}

	for _, tc := range testCases {
		gotNecessaryPacks := getMinNecessaryPacks(tc.Items)
		if !reflect.DeepEqual(tc.WantNecessaryPacks, gotNecessaryPacks) {
			t.Fatalf("For %v items, expected: %v, got %v", tc.Items, tc.WantNecessaryPacks, gotNecessaryPacks)
		}
	}
}
