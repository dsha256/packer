package packer

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sort"
)

const (
	ErrorNegativeOrZeroSize = "size must be more than 0"
	ErrorDuplicatedSizes    = "size already exists or incoming sizes contains duplications"
	ErrorZeroSizesQuantity  = "sizes must be more than 0 in quantity"
	ErrorSizeDoesNotExist   = "size does not exist"
)

var SortedSizes = []int{250, 500, 1000, 2000, 5000}

type SizerService struct {
	SortedSizes []int
}

func NewSizer(sizes []int) *SizerService {
	if !slices.IsSorted(sizes) {
		slices.Sort(sizes)
	}

	sizesSrvc := &SizerService{
		SortedSizes: sizes,
	}
	sort.Ints(sizesSrvc.SortedSizes)

	return sizesSrvc
}

func (sizes *SizerService) ListSizes() []int {
	return sizes.SortedSizes
}

func (sizes *SizerService) AddSize(ctx context.Context, sizeToAdd int) ([]int, error) {
	if sizeToAdd <= 0 {
		return []int{}, errors.New(ErrorNegativeOrZeroSize)
	}
	if sizes.Exists(sizeToAdd) {
		slog.ErrorContext(ctx,
			ErrorDuplicatedSizes,
			slog.Any("incoming_size", sizeToAdd),
			slog.Any("existing_sizes", sizes.SortedSizes),
		)
		return []int{}, errors.New(ErrorDuplicatedSizes)
	}

	sizes.SortedSizes = insertSorted(sizes.SortedSizes, sizeToAdd)

	return sizes.SortedSizes, nil
}

func (sizes *SizerService) PutSizes(ctx context.Context, sizesToPut []int) ([]int, error) {
	if len(sizesToPut) == 0 {
		return []int{}, errors.New(ErrorZeroSizesQuantity)
	}

	sizesWeights := make(map[int]int)
	for _, size := range sizesToPut {
		if size <= 0 {
			slog.ErrorContext(ctx,
				ErrorNegativeOrZeroSize,
				slog.Any("incoming_size", size),
				slog.Any("existing_sizes", sizes.SortedSizes),
			)
			return []int{}, errors.New(ErrorNegativeOrZeroSize)
		}
		if _, exists := sizesWeights[size]; exists {
			slog.ErrorContext(ctx,
				ErrorDuplicatedSizes,
				slog.Any("incoming_size", size),
				slog.Any("existing_sizes", sizes.SortedSizes),
			)
			return []int{}, errors.New(ErrorDuplicatedSizes)
		}
		sizesWeights[size] = 1
	}

	sizes.SortedSizes = []int{}
	slices.Sort(sizesToPut)
	sizes.SortedSizes = append(sizes.SortedSizes, sizesToPut...)

	return sizes.SortedSizes, nil
}

func (sizes *SizerService) DeleteSize(ctx context.Context, sizeToDelete int) ([]int, error) {
	if sizeToDelete <= 0 {
		slog.ErrorContext(ctx,
			ErrorNegativeOrZeroSize,
			slog.Any("incoming_size", sizeToDelete),
			slog.Any("existing_sizes", sizes.SortedSizes),
		)
		return []int{}, errors.New(ErrorNegativeOrZeroSize)
	}

	if !sizes.Exists(sizeToDelete) {
		slog.ErrorContext(ctx,
			ErrorSizeDoesNotExist,
			slog.Any("incoming_size", sizeToDelete),
			slog.Any("existing_sizes", sizes.SortedSizes),
		)
		return []int{}, errors.New(ErrorSizeDoesNotExist)
	}

	indexOfSizeToDelete, _ := slices.BinarySearch(sizes.SortedSizes, sizeToDelete)
	sizes.SortedSizes = slices.Delete(sizes.SortedSizes, indexOfSizeToDelete, indexOfSizeToDelete+1)

	return sizes.SortedSizes, nil
}

func (sizes *SizerService) Exists(sizeToCheckFor int) bool {
	_, exists := slices.BinarySearch(sizes.SortedSizes, sizeToCheckFor)
	return exists
}

// insertSorted inserts element to the slice in a sorted passion.
func insertSorted(targetSlice []int, element int) []int {
	i := sort.Search(len(targetSlice), func(i int) bool { return targetSlice[i] > element })
	targetSlice = append(targetSlice, 0)
	copy(targetSlice[i+1:], targetSlice[i:])
	targetSlice[i] = element
	return targetSlice
}
