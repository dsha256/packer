package packer

import (
	"context"

	"github.com/dsha256/packer/internal/types"
)

type Packer interface {
	ListPacketSizes(ctx context.Context) ([]types.PacketSize, error)
	SetPacketSizes(ctx context.Context, sizes []types.PacketSize) error
	GetOptimalPackets(ctx context.Context, items int) (map[types.PacketSize]types.PacketQuantity, error)
}
