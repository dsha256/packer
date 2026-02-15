package packer

import (
	"context"
	"sort"
	"sync"

	"github.com/dsha256/packer/internal/types"
)

type packer struct {
	packetSizes     []types.PacketSize
	packetSizesLock sync.Mutex
}

func New() Packer {
	return &packer{
		packetSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
	}
}

func (s *packer) ListPacketSizes(_ context.Context) ([]types.PacketSize, error) {
	return s.packetSizes, nil
}

func (s *packer) SetPacketSizes(_ context.Context, sizes []types.PacketSize) error {
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] < sizes[j]
	})

	s.packetSizesLock.Lock()
	s.packetSizes = sizes
	s.packetSizesLock.Unlock()

	return nil
}

func (s *packer) GetOptimalPackets(_ context.Context, items int) (map[types.PacketSize]types.PacketQuantity, error) {
	return CalculateOptimalPacketsForItemsV1(&CalculateOptimalPacketsForItemsParams{
		Items:       items,
		PacketSizes: s.packetSizes,
	}), nil
}
