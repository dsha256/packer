# Benchmarking Algorithms: `CalculateOptimalPacketsForItems`

This document presents the benchmarking results and analysis of two algorithms, `CalculateOptimalPacketsForItemsV1` and `CalculateOptimalPacketsForItemsV2`, used to calculate the optimal allocation of packets for various item quantities. The benchmarking was conducted to evaluate their performance in terms of execution time, memory allocation, and scalability across different input sizes.

## Benchmarking Summary

### Command Used for Benchmarking
The following command was used to perform the benchmark tests:

```bash
go test -v internal/packer/*.go -bench=. -run=xxx -benchmem -benchtime=5s -count=5
```

### Test Environment
- **Operating System:** macOS Sonoma
- **Processor Architecture:** Apple Silicon (arm64)
- **CPU Model:** Apple M3 Pro
- **Go Version:** Go SDK 1.24.2, devel

### Benchmark Code
The benchmarking was conducted using the following Go code:

```go
package packer

import (
	"testing"
	"github.com/dsha256/packer/internal/types"
)

func benchmarkCalculateOptimalPacketsForItems(b *testing.B, calculateFunc func(params *CalculateOptimalPacketsForItemsParams) map[types.PacketSize]types.PacketQuantity) {
	testCases := []struct {
		Name        string
		Items       int
		PacketSizes []types.PacketSize
	}{
		{"ExtraSmallQuantityOfItems", 250, []types.PacketSize{250, 500, 1000, 2000, 5000}},
		{"SmallQuantityOfItems", 12001, []types.PacketSize{250, 500, 1000, 2000, 5000}},
		{"MediumQuantityOfItems", 100_000, []types.PacketSize{23, 31, 53}},
		{"LargeQuantityOfItems", 500_000, []types.PacketSize{23, 31, 53}},
		{"ExtraLargeQuantityOfItems", 10_000_000, []types.PacketSize{250, 500, 1000, 2000, 5000}},
	}

	for _, testCase := range testCases {
		b.Run(testCase.Name, func(b *testing.B) {
			params := &CalculateOptimalPacketsForItemsParams{
				Items:       testCase.Items,
				PacketSizes: testCase.PacketSizes,
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				calculateFunc(params)
			}
		})
	}
}

func Benchmark_CalculateOptimalPacketsForItemsV2(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItems(b, CalculateOptimalPacketsForItemsV2)
}

func Benchmark_CalculateOptimalPacketsForItemsV1(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItems(b, CalculateOptimalPacketsForItemsV1)
}
```

---

## Benchmark Results

The detailed benchmark results for each algorithm and input size are as follows:

### **`CalculateOptimalPacketsForItemsV1`**

| Input Size                         | Iterations | Execution Time (ns/op) | Memory Allocated (B/op) | Allocations (allocs/op) |
|------------------------------------|------------|-------------------------|--------------------------|-------------------------|
| **ExtraSmallQuantityOfItems**      | ≈ 242,000  | ≈ 25,000               | ≈ 98,500                | 5                       |
| **SmallQuantityOfItems**           | ≈ 83,000   | ≈ 72,400               | ≈ 278,800               | 5                       |
| **MediumQuantityOfItems**          | 10,000     | ≈ 581,000              | ≈ 1,605,800             | 5                       |
| **LargeQuantityOfItems**           | 2,200      | ≈ 2,740,000            | ≈ 8,012,000             | 5                       |
| **ExtraLargeQuantityOfItems**      | 160        | ≈ 37,170,000           | ≈ 160,088,300           | 5                       |

---

### **`CalculateOptimalPacketsForItemsV2`**

| Input Size                         | Iterations   | Execution Time (ns/op) | Memory Allocated (B/op) | Allocations (allocs/op) |
|------------------------------------|--------------|-------------------------|--------------------------|-------------------------|
| **ExtraSmallQuantityOfItems**      | ≈ 18,300,000 | ≈ 330                 | ≈ 584                   | 15                      |
| **SmallQuantityOfItems**           | ≈ 614,000    | ≈ 9,720               | ≈ 14,024                | 146                     |
| **MediumQuantityOfItems**          | ≈ 340        | ≈ 17,400,000          | ≈ 14,916,430            | 200,789                 |
| **LargeQuantityOfItems**           | ≈ 46         | ≈ 117,780,000         | ≈ 109,746,000           | 1,007,962               |
| **ExtraLargeQuantityOfItems**      | ≈ 738        | ≈ 8,120,000           | ≈ 7,140,622             | 80,576                  |

---

## Insights and Analysis

### General Observations
1. **Execution Time:**
   - `CalculateOptimalPacketsForItemsV2` demonstrated significantly better performance for smaller inputs (`ExtraSmallQuantityOfItems` and `SmallQuantityOfItems`), with execution times reduced by multiple orders of magnitude compared to `CalculateOptimalPacketsForItemsV1`.
   - However, for larger inputs (`MediumQuantityOfItems` and above), `CalculateOptimalPacketsForItemsV2` took longer to execute per operation compared to its counterpart.

2. **Memory Usage:**
   - `CalculateOptimalPacketsForItemsV2` utilized orders of magnitude less memory compared to `CalculateOptimalPacketsForItemsV1` for the `ExtraSmallQuantityOfItems` and `SmallQuantityOfItems` scenarios.
   - However, for medium to large datasets, the memory usage and allocations grew significantly. This indicates that while `V2` performs better in simpler cases, it requires optimization for large-scale scenarios.

3. **Allocations Per Operation:**
   - `CalculateOptimalPacketsForItemsV2` made significantly more allocations per operation compared to `CalculateOptimalPacketsForItemsV1`. This may be indicative of its design relying more heavily on dynamic memory allocation, which could impact performance on larger datasets.

---

### Key Considerations
1. **Algorithm Design:**
   - While `CalculateOptimalPacketsForItemsV2` is optimized for smaller input sizes, its design should be revisited to improve performance and scalability for medium-to-large datasets.
   - Profiling for memory and allocation bottlenecks in `V2` could yield key improvements.

2. **Use Case Scenarios:**
   - For systems where inputs are consistently small, `V2` is the obvious choice due to its superior performance and lower memory overhead.
   - For high-load or large dataset scenarios, `V1` may provide better scalability under the current implementation.

3. **Potential Improvements:**
   - Explore more efficient data structures or caching mechanisms in `V2` to reduce allocations and improve performance for larger datasets.
   - Consider hybrid approaches that switch between `V1` and `V2` based on input size.

---

### Recommendations
- **Scenarios with Small Inputs:** Use `CalculateOptimalPacketsForItemsV2`. Its performance and memory efficiency in these cases make it the ideal choice.
- **Scenarios with Medium to Large Inputs:** Default to `CalculateOptimalPacketsForItemsV1` unless `V2` can be optimized to close the gap in scalability.

---

## Conclusion

The benchmarking results provide valuable insights into the performance characteristics of both algorithms. `CalculateOptimalPacketsForItemsV2` shows great promise for small-scale operations due to its exceptional performance. However, improvements are required to make it competitive in larger-scale use cases. Future iterations should aim to balance the strengths of both algorithms.