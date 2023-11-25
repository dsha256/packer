package packer

import "context"

type AggregatePacker interface {
	Sizer
	Packer
}

type Aggregator struct {
	Sizer
	Packer
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		Sizer:  NewSizerService(SortedSizes),
		Packer: NewPacketsService(),
	}
}

type Sizer interface {
	ListSizes() []int
	AddSize(ctx context.Context, sizeToAdd int) ([]int, error)
	PutSizes(ctx context.Context, sizesToPut []int) ([]int, error)
	DeleteSize(ctx context.Context, sizeToDelete int) ([]int, error)
	Exists(sizeToCheckFor int) bool
}

type Packer interface {
	GetPackets(ctx context.Context, itemsToPack int) (map[int]int, error)
}
