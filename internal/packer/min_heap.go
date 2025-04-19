package packer

type HeapElement struct {
	total    int
	numPacks int
}

type MinHeap []HeapElement

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i].total < h[j].total || (h[i].total == h[j].total && h[i].numPacks < h[j].numPacks)
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x any) {
	element, ok := x.(HeapElement) // Separate the type assertion
	if !ok {
		return
	}
	*h = append(*h, element)
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
