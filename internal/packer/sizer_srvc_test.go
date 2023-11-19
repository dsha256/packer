package packer

import (
	"context"
	"log"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func newSizer(sizes []int) *SizerService {
	return NewSizerService(sizes)
}

func Test_NewSizer(t *testing.T) {
	sizer := newSizer(SortedSizes)

	require.True(t, len(sizer.SortedSizes) > 0)
	require.True(t, slices.Equal(sizer.SortedSizes, SortedSizes))
	require.True(t, slices.IsSorted(sizer.SortedSizes))

	sizer = newSizer([]int{1, 2, 10, 7})
	require.True(t, slices.IsSorted(sizer.SortedSizes))
}

func TestSizerService_ListSizes(t *testing.T) {
	sizer := newSizer(SortedSizes)
	sizes := sizer.ListSizes()
	require.True(t, len(sizes) > 0)
	require.True(t, slices.IsSorted(sizes))
}

func TestSizerService_AddSize(t *testing.T) {
	testCases := []struct {
		name               string
		initialSortedSizes []int
		size               int
		wantErr            bool
		checkResult        func(t *testing.T, size int, sortedSizes []int, err error)
	}{
		{
			name:               "OK add in the end",
			initialSortedSizes: []int{250, 500, 1000, 2000, 5000},
			size:               10000,
			wantErr:            false,
			checkResult: func(t *testing.T, size int, sortedSizes []int, err error) {
				require.True(t, slices.Equal(sortedSizes, insertSorted(SortedSizes, size)))
			},
		},
		{
			name:               "OK add in the middle",
			initialSortedSizes: []int{250, 500, 1000, 2000, 5000},
			size:               1500,
			wantErr:            false,
			checkResult: func(t *testing.T, size int, sortedSizes []int, err error) {
				require.True(t, slices.Equal(sortedSizes, insertSorted(SortedSizes, size)))
			},
		},
		{
			name:               "ERR add zero",
			initialSortedSizes: []int{250, 500, 1000, 2000, 5000},
			size:               0,
			wantErr:            true,
			checkResult: func(t *testing.T, size int, sortedSizes []int, err error) {
				require.Error(t, err)
				require.Equal(t, err.Error(), ErrorNegativeOrZeroSize)
			},
		},
		{
			name:               "ERR add negative size",
			initialSortedSizes: []int{250, 500, 1000, 2000, 5000},
			size:               -1,
			wantErr:            true,
			checkResult: func(t *testing.T, size int, sortedSizes []int, err error) {
				require.Error(t, err)
				require.Equal(t, err.Error(), ErrorNegativeOrZeroSize)
			},
		},
		{
			name:               "ERR add existing size",
			initialSortedSizes: []int{250, 500, 1000, 2000, 5000},
			size:               500,
			wantErr:            true,
			checkResult: func(t *testing.T, size int, sortedSizes []int, err error) {
				require.Error(t, err)
				require.Equal(t, err.Error(), ErrorDuplicatedSizes)
			},
		},
	}

	for index := range testCases {
		tc := testCases[index]
		sizer := newSizer(tc.initialSortedSizes)
		t.Run(tc.name, func(t *testing.T) {
			sortedSizes, err := sizer.AddSize(context.Background(), tc.size)
			if tc.wantErr && err == nil {
				log.Println("wanted err but got nil instead")
				t.Fail()
			}
			if !tc.wantErr && err != nil {
				log.Println("wanted nil but got err instead")
			}
			tc.checkResult(t, tc.size, sortedSizes, err)
		})
	}
}

func TestSizerService_PutSizes(t *testing.T) {
	testCases := []struct {
		name        string
		sizesToPut  []int
		checkResult func(t *testing.T, sizesToPut []int, sortedSizes []int, err error)
	}{
		{
			name:       "OK put acceptable sorted sizes (ascending)",
			sizesToPut: []int{1, 10, 100, 200},
			checkResult: func(t *testing.T, sizesToPut []int, sortedSizes []int, err error) {
				require.NoError(t, err)
				require.True(t, len(sizesToPut) == len(sortedSizes))
				require.True(t, slices.IsSorted(sortedSizes))
				require.True(t, slices.Equal(sizesToPut, sortedSizes))
			},
		},
		{
			name:       "OK put acceptable sorted sizes (descending)",
			sizesToPut: []int{200, 100, 10, 1},
			checkResult: func(t *testing.T, sizesToPut []int, sortedSizes []int, err error) {
				require.NoError(t, err)
				require.True(t, len(sizesToPut) == len(sortedSizes))
				require.True(t, slices.IsSorted(sortedSizes))
				require.True(t, slices.Equal(sizesToPut, sortedSizes))
			},
		},
		{
			name:       "OK put acceptable not sorted sizes",
			sizesToPut: []int{200, 1000, 10, 1},
			checkResult: func(t *testing.T, sizesToPut []int, sortedSizes []int, err error) {
				require.NoError(t, err)
				require.True(t, len(sizesToPut) == len(sortedSizes))
				require.True(t, slices.IsSorted(sortedSizes))
				require.True(t, slices.Equal(sizesToPut, sortedSizes))
			},
		},
		{
			name:       "ERR put empty sizes",
			sizesToPut: []int{},
			checkResult: func(t *testing.T, sizesToPut []int, sortedSizes []int, err error) {
				require.Error(t, err)
				require.Equal(t, err.Error(), ErrorZeroSizesQuantity)
				require.True(t, len(sortedSizes) == 0)
			},
		},
		{
			name:       "ERR put non acceptable sizes - 0",
			sizesToPut: []int{1, 3, 0},
			checkResult: func(t *testing.T, sizesToPut []int, sortedSizes []int, err error) {
				require.Error(t, err)
				require.Equal(t, err.Error(), ErrorNegativeOrZeroSize)
				require.True(t, len(sortedSizes) == 0)
			},
		},
		{
			name:       "ERR put non acceptable sizes - (-1)",
			sizesToPut: []int{1, 3, -1},
			checkResult: func(t *testing.T, sizesToPut []int, sortedSizes []int, err error) {
				require.Error(t, err)
				require.Equal(t, err.Error(), ErrorNegativeOrZeroSize)
				require.True(t, len(sortedSizes) == 0)
			},
		},
		{
			name:       "ERR put duplicated sizes",
			sizesToPut: []int{1, 3, 1},
			checkResult: func(t *testing.T, sizesToPut []int, sortedSizes []int, err error) {
				require.Error(t, err)
				require.Equal(t, err.Error(), ErrorDuplicatedSizes)
				require.True(t, len(sortedSizes) == 0)
			},
		},
	}

	for index := range testCases {
		sizer := newSizer(SortedSizes)
		tc := testCases[index]
		t.Run(tc.name, func(t *testing.T) {
			sortedSizes, err := sizer.PutSizes(context.Background(), tc.sizesToPut)
			tc.checkResult(t, tc.sizesToPut, sortedSizes, err)
		})
	}
}

func TestSizerService_DeleteSize(t *testing.T) {
	testCases := []struct {
		name               string
		initialSortedSizes []int
		sizeToDelete       int
		checkResult        func(t *testing.T, sizesToDelete int, sortedSizes []int, err error)
	}{
		{
			name:               "OK delete existing size",
			initialSortedSizes: []int{1, 2, 3, 10},
			sizeToDelete:       3,
			checkResult: func(t *testing.T, sizesToDelete int, sortedSizes []int, err error) {
				require.NoError(t, err)
				require.False(t, slices.Contains(sortedSizes, sizesToDelete))
			},
		},
		{
			name:               "ERR delete non existing size",
			initialSortedSizes: []int{1, 2, 3, 10},
			sizeToDelete:       7,
			checkResult: func(t *testing.T, sizesToDelete int, sortedSizes []int, err error) {
				require.Error(t, err)
				require.Equal(t, err.Error(), ErrorSizeDoesNotExist)
			},
		},
		{
			name:               "ERR delete negative size",
			initialSortedSizes: []int{1, 2, 3, 10},
			sizeToDelete:       -1,
			checkResult: func(t *testing.T, sizesToDelete int, sortedSizes []int, err error) {
				require.Error(t, err)
				require.Equal(t, err.Error(), ErrorNegativeOrZeroSize)
			},
		},
	}

	for index := range testCases {
		tc := testCases[index]
		sizer := newSizer(tc.initialSortedSizes)
		t.Run(tc.name, func(t *testing.T) {
			sortedSizes, err := sizer.DeleteSize(context.Background(), tc.sizeToDelete)
			tc.checkResult(t, tc.sizeToDelete, sortedSizes, err)
		})
	}
}
