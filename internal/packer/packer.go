package packer

import "context"

type Sizer interface {
	AddSize(ctx context.Context, sizeToAdd int) ([]int, error)
	PutSizes(ctx context.Context, sizesToPut []int) ([]int, error)
	DeleteSize(ctx context.Context, sizeToDelete int) ([]int, error)
	Exists(sizeToCheckFor int) bool
}

type Packer interface {
	GetPackets(ctx context.Context, itemsToPack int) (map[int]int, error)
}
