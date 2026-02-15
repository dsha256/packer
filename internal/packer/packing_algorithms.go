package packer

import (
	"container/heap"
	"math"

	"github.com/dsha256/packer/internal/types"
)

type CalculateOptimalPacketsForItemsParams struct {
	PacketSizes []types.PacketSize
	Items       int
}

func CalculateOptimalPacketsForItemsV1(params *CalculateOptimalPacketsForItemsParams) map[types.PacketSize]types.PacketQuantity {
	result := make(map[types.PacketSize]types.PacketQuantity)

	sizes := make([]int, len(params.PacketSizes))
	for i, ps := range params.PacketSizes {
		sizes[i] = int(ps)
	}

	maxSize := sizes[len(sizes)-1]
	maxSum := params.Items + maxSize

	dpPacks := make([]int, maxSum+1)
	prevPacket := make([]int, maxSum+1)
	for i := range dpPacks {
		dpPacks[i] = math.MaxInt32
	}
	dpPacks[0] = 0

	for s := 1; s <= maxSum; s++ {
		for _, sz := range sizes {
			if s >= sz && dpPacks[s-sz] != math.MaxInt32 {
				if cand := dpPacks[s-sz] + 1; cand < dpPacks[s] {
					dpPacks[s] = cand
					prevPacket[s] = sz
				}
			}
		}
	}

	bestSum := -1
	for s := params.Items; s <= maxSum; s++ {
		if dpPacks[s] < math.MaxInt32 {
			bestSum = s

			break
		}
	}
	if bestSum < 0 {
		result[types.PacketSize(sizes[0])] = 1

		return result
	}

	countMap := make(map[int]int)
	for cur := bestSum; cur > 0; cur -= prevPacket[cur] {
		countMap[prevPacket[cur]]++
	}

	for _, sz := range sizes {
		if cnt, ok := countMap[sz]; ok && cnt > 0 {
			result[types.PacketSize(sz)] = types.PacketQuantity(cnt)
		}
	}

	return result
}

func CalculateOptimalPacketsForItemsV2(params *CalculateOptimalPacketsForItemsParams) map[types.PacketSize]types.PacketQuantity {
	items := params.Items
	sizes := params.PacketSizes

	minNumPacks := make(map[int]int)
	predecessor := make(map[int]struct {
		prevT int
		s     types.PacketSize
	})

	minHeap := &MinHeap{}
	heap.Init(minHeap)
	heap.Push(minHeap, HeapElement{0, 0})
	minNumPacks[0] = 0

	for minHeap.Len() > 0 {
		popped := heap.Pop(minHeap)
		heapElement, ok := popped.(HeapElement)
		if !ok {
			return make(map[types.PacketSize]types.PacketQuantity)
		}
		total := heapElement.total
		numPacks := heapElement.numPacks

		if numPacks > minNumPacks[total] {
			continue
		}

		if total >= items {
			result := make(map[types.PacketSize]types.PacketQuantity)
			currentTotal := total
			for currentTotal != 0 {
				pred := predecessor[currentTotal]
				s := pred.s
				result[s]++
				currentTotal = pred.prevT
			}

			return result
		}

		for _, size := range sizes {
			newTotal := total + int(size)
			newNumPacks := numPacks + 1
			if _, ok := minNumPacks[newTotal]; !ok || newNumPacks < minNumPacks[newTotal] {
				minNumPacks[newTotal] = newNumPacks
				predecessor[newTotal] = struct {
					prevT int
					s     types.PacketSize
				}{total, size}
				heap.Push(minHeap, HeapElement{newTotal, newNumPacks})
			}
		}
	}

	return make(map[types.PacketSize]types.PacketQuantity)
}
