# Benchmarking Algorithms: `CalculateOptimalPacketsForItems`

This document presents the benchmarking results and analysis of two algorithms, `CalculateOptimalPacketsForItemsV1` and `CalculateOptimalPacketsForItemsV2`, used to calculate the optimal allocation of packets for various item quantities. The benchmarking was conducted to evaluate their performance in terms of execution time, memory allocation, and scalability across different input sizes and packet size configurations.

---

### Summary of Complexities

| **Algorithm**                     | **Time Complexity**                                   | **Space Complexity**       | **Approach**                           |
|------------------------------------|------------------------------------------------------|-----------------------------|----------------------------------------|
| **`CalculateOptimalPacketsForItemsV1`** | `O((items + maxPacketSize) * len(packetSizes))`      | `O(items + maxPacketSize)` | Dynamic Programming (Backtracking)    |
| **`CalculateOptimalPacketsForItemsV2`** | `O((items + maxPacketSize) * len(packetSizes) * log(items + maxPacketSize))` | `O(items + maxPacketSize)` | Dijkstra's Algorithm with Min-Heap    |

---

### Key Differences in Approach

| **Factor**               | **V1 (DP)**                                | **V2 (Min-Heap / Dijkstra)**             |
|--------------------------|--------------------------------------------|------------------------------------------|
| **Algorithm Type**       | Dynamic Programming                        | Priority Queue + Greedy Traversal (Dijkstra) |
| **Main Data Structure**  | Arrays (`dpPacks`, `prevPacket`)            | Heap (`MinHeap`) + Maps (`minNumPacks`)  |
| **Backtracking**         | Uses `prevPacket` to reconstruct solution. | Uses `predecessor` map to reconstruct solution. |
| **Efficiency**           | Processes all totals up to `maxSum`.       | Prioritizes smaller totals with fewer packets first. |

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
- **Go Version:** Go SDK 1.25, devel

### Benchmark Code
The benchmarking was conducted using the following Go code:

```go
package packer_test

import (
	"testing"

	"github.com/dsha256/packer/internal/packer"
	"github.com/dsha256/packer/internal/types"
)

func benchmarkCalculateOptimalPacketsForItemsWithProductOfTenSizes(b *testing.B, calculateFunc func(params *packer.CalculateOptimalPacketsForItemsParams) map[types.PacketSize]types.PacketQuantity) {
	b.Helper()

	testCases := []struct {
		Name        string
		PacketSizes []types.PacketSize
		Items       int
	}{
		{
			Name:        "~50k_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       50_123,
		},
		{
			Name:        "~100k_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       100_123,
		},
		{
			Name:        "~250k_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       250_123,
		},
		{
			Name:        "~500k_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       500_123,
		},
		{
			Name:        "~1M_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       1_123_123,
		},
		{
			Name:        "~10M_Items",
			PacketSizes: []types.PacketSize{250, 500, 1000, 2000, 5000},
			Items:       10_123_123,
		},
	}

	for _, testCase := range testCases {
		b.Run(testCase.Name, func(b *testing.B) {
			params := &packer.CalculateOptimalPacketsForItemsParams{
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

func benchmarkCalculateOptimalPacketsForItemsWithPrimeSizes(b *testing.B, calculateFunc func(params *packer.CalculateOptimalPacketsForItemsParams) map[types.PacketSize]types.PacketQuantity) {
	b.Helper()

	testCases := []struct {
		Name        string
		PacketSizes []types.PacketSize
		Items       int
	}{
		{
			Name:        "~50k_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       50_123,
		},
		{
			Name:        "~100k_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       100_123,
		},
		{
			Name:        "~250k_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       250_123,
		},
		{
			Name:        "~500k_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       500_123,
		},
		{
			Name:        "~1M_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       1_123_123,
		},
		{
			Name:        "~10M_Items",
			PacketSizes: []types.PacketSize{251, 503, 997, 2003, 4999},
			Items:       10_123_123,
		},
	}

	for _, testCase := range testCases {
		b.Run(testCase.Name, func(b *testing.B) {
			params := &packer.CalculateOptimalPacketsForItemsParams{
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

func Benchmark_CalculateOptimalPacketsForItemsV2_ProductOfTenSizes(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItemsWithProductOfTenSizes(b, packer.CalculateOptimalPacketsForItemsV2)
}

func Benchmark_CalculateOptimalPacketsForItemsV1_ProductOfTenSizes(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItemsWithProductOfTenSizes(b, packer.CalculateOptimalPacketsForItemsV1)
}

func Benchmark_CalculateOptimalPacketsForItemsV2_PrimeSizes(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItemsWithPrimeSizes(b, packer.CalculateOptimalPacketsForItemsV2)
}

func Benchmark_CalculateOptimalPacketsForItemsV1_PrimeSizes(b *testing.B) {
	benchmarkCalculateOptimalPacketsForItemsWithPrimeSizes(b, packer.CalculateOptimalPacketsForItemsV1)
}
```

---

## Benchmark Results

The detailed benchmark results for each algorithm and input size are as follows:

### Product of Ten Sizes (250, 500, 1000, 2000, 5000)

#### V2 Algorithm

| Input Size | Iterations | Execution Time (ns/op) | Memory Allocated (B/op) | Allocations (allocs/op) |
|------------|------------|------------------------|-------------------------|------------------------|
| ~50k Items | ≈ 170,000  | ≈ 37,000              | ≈ 30,344               | 454                    |
| ~100k Items| ≈ 70,000   | ≈ 84,000              | ≈ 59,848               | 858                    |
| ~250k Items| ≈ 25,000   | ≈ 240,000             | ≈ 216,361              | 2,072                  |
| ~500k Items| ≈ 12,000   | ≈ 475,000             | ≈ 431,465              | 4,090                  |
| ~1M Items  | ≈ 5,700    | ≈ 1,050,000           | ≈ 877,419              | 9,108                  |
| ~10M Items | ≈ 600      | ≈ 9,600,000           | ≈ 7,156,376            | 81,562                 |

#### V1 Algorithm

| Input Size | Iterations | Execution Time (ns/op) | Memory Allocated (B/op) | Allocations (allocs/op) |
|------------|------------|------------------------|-------------------------|------------------------|
| ~50k Items | ≈ 22,000   | ≈ 265,000             | ≈ 884,978              | 5                      |
| ~100k Items| ≈ 10,000   | ≈ 500,000             | ≈ 1,687,796            | 5                      |
| ~250k Items| ≈ 4,800    | ≈ 1,190,000           | ≈ 4,096,250            | 5                      |
| ~500k Items| ≈ 2,800    | ≈ 2,150,000           | ≈ 8,093,948            | 5                      |
| ~1M Items  | ≈ 1,300    | ≈ 4,650,000           | ≈ 18,055,425           | 5                      |
| ~10M Items | ≈ 150      | ≈ 40,000,000          | ≈ 162,054,404          | 5                      |

### Prime Number Sizes (251, 503, 997, 2003, 4999)

#### V2 Algorithm

| Input Size | Iterations | Execution Time (ns/op) | Memory Allocated (B/op) | Allocations (allocs/op) |
|------------|------------|------------------------|-------------------------|------------------------|
| ~50k Items | ≈ 470      | ≈ 12,800,000          | ≈ 7,591,019            | 85,938                 |
| ~100k Items| ≈ 200      | ≈ 30,000,000          | ≈ 15,115,240           | 190,490                |
| ~250k Items| ≈ 57       | ≈ 105,000,000         | ≈ 55,073,280           | 493,567                |
| ~500k Items| ≈ 22       | ≈ 252,000,000         | ≈ 109,810,000          | 997,653                |
| ~1M Items  | ≈ 9        | ≈ 625,000,000         | ≈ 223,639,832          | 2,251,859              |
| ~10M Items | ≈ 1        | ≈ 7,500,000,000       | ≈ 1,824,129,464        | 20,366,553             |

#### V1 Algorithm

| Input Size | Iterations | Execution Time (ns/op) | Memory Allocated (B/op) | Allocations (allocs/op) |
|------------|------------|------------------------|-------------------------|------------------------|
| ~50k Items | ≈ 15,000   | ≈ 400,000             | ≈ 884,979              | 5                      |
| ~100k Items| ≈ 7,700    | ≈ 780,000             | ≈ 1,687,799            | 5                      |
| ~250k Items| ≈ 3,100    | ≈ 1,900,000           | ≈ 4,096,255            | 5                      |
| ~500k Items| ≈ 1,600    | ≈ 3,600,000           | ≈ 8,093,955            | 5                      |
| ~1M Items  | ≈ 1,300    | ≈ 4,650,000           | ≈ 18,055,426           | 5                      |
| ~10M Items | ≈ 150      | ≈ 40,000,000          | ≈ 162,054,404          | 5                      |

---

## Insights and Analysis

### General Observations

1. **Execution Time:**
   - With product-of-ten packet sizes, V2 performs better for inputs up to ~1M items, after which V1 becomes more efficient.
   - With prime number packet sizes, V2 performs poorly across all input sizes, with execution times orders of magnitude higher than V1.
   - V1's performance is relatively consistent regardless of packet size configuration, making it more robust for different use cases.

2. **Memory Usage:**
   - With product-of-ten packet sizes, V2 uses less memory than V1 for inputs up to ~500k items.
   - With prime number packet sizes, V2 uses significantly more memory than V1 across all input sizes.
   - V1's memory usage scales linearly with input size, while V2's memory usage grows more rapidly, especially with prime number packet sizes.

3. **Allocations Per Operation:**
   - V1 consistently makes exactly 5 allocations per operation, regardless of input size or packet size configuration.
   - V2's allocation count varies widely based on input size and packet size configuration, ranging from 454 to over 20 million allocations.
   - With prime number packet sizes, V2's allocation count grows dramatically with input size.

4. **Packet Size Configuration Impact:**
   - The choice of packet sizes has a dramatic impact on V2's performance, with prime number sizes causing severe performance degradation.
   - V1's performance is relatively consistent regardless of packet size configuration, making it more robust for different use cases.

---

### Key Considerations

1. **Algorithm Design:**
   - While `CalculateOptimalPacketsForItemsV2` is optimized for smaller input sizes and product-of-ten packet sizes, its design should be revisited to improve performance and scalability for medium-to-large datasets and prime number packet sizes.
   - Profiling for memory and allocation bottlenecks in `V2` could yield key improvements.

2. **Use Case Scenarios:**
   - For systems where inputs are consistently small and packet sizes are products of ten, `V2` is the obvious choice due to its superior performance and lower memory overhead.
   - For high-load or large dataset scenarios, or when using prime number packet sizes, `V1` provides better scalability under the current implementation.

3. **Potential Improvements:**
   - Explore more efficient data structures or caching mechanisms in `V2` to reduce allocations and improve performance for larger datasets and prime number packet sizes.
   - Consider hybrid approaches that switch between `V1` and `V2` based on input size and packet size configuration.

---

## Hybrid Approach Recommendation

Based on the comprehensive benchmark results, we recommend implementing a hybrid approach that automatically selects the most appropriate algorithm based on input characteristics:

```go
func CalculateOptimalPackets(items int, packetSizes []PacketSize) map[PacketSize]PacketQuantity {
    // Check if packet sizes are prime numbers
    isPrimeSizes := true
    for _, size := range packetSizes {
        if !isPrime(size) {
            isPrimeSizes = false
            break
        }
    }
    
    // Use V2 for small to medium datasets with product-of-ten sizes
    if (items < 100_000 || (items < 1_000_000 && !isPrimeSizes)) {
        return CalculateOptimalPacketsForItemsV2(&CalculateOptimalPacketsForItemsParams{
            Items: items,
            PacketSizes: packetSizes,
        })
    }
    
    // Use V1 for large datasets or prime number sizes
    return CalculateOptimalPacketsForItemsV1(&CalculateOptimalPacketsForItemsParams{
        Items: items,
        PacketSizes: packetSizes,
    })
}

// Helper function to check if a number is prime
func isPrime(n int) bool {
    if n <= 1 {
        return false
    }
    if n <= 3 {
        return true
    }
    if n%2 == 0 || n%3 == 0 {
        return false
    }
    
    for i := 5; i*i <= n; i += 6 {
        if n%i == 0 || n%(i+2) == 0 {
            return false
        }
    }
    return true
}
```

This hybrid approach provides the following benefits:

1. **Optimal Performance for Small Inputs**: Uses V2 for small to medium datasets with product-of-ten packet sizes, where it performs best.
2. **Robust Performance for Large Inputs**: Uses V1 for large datasets or when using prime number packet sizes, where it performs better.
3. **Adaptive Selection**: Automatically selects the best algorithm based on input characteristics, without requiring manual intervention.
4. **Future-Proofing**: Can be easily updated as algorithms are improved or new algorithms are added.

---

## Conclusion

The benchmarking results provide valuable insights into the performance characteristics of both algorithms across different input sizes and packet size configurations. `CalculateOptimalPacketsForItemsV2` shows great promise for small to medium-scale operations with product-of-ten packet sizes, while `CalculateOptimalPacketsForItemsV1` provides more consistent and robust performance for large-scale operations and prime number packet sizes.

The hybrid approach recommended in this document offers a practical solution that leverages the strengths of both algorithms, ensuring optimal performance across all use cases. By automatically selecting the most appropriate algorithm based on input characteristics, we can achieve the best of both worlds: the efficiency of V2 for small to medium inputs with product-of-ten packet sizes, and the robustness of V1 for large inputs and prime number packet sizes.

Future development efforts should focus on optimizing V2 for large inputs and prime number packet sizes, potentially incorporating techniques from V1 to improve its performance in these scenarios. Additionally, further benchmarking with real-world usage patterns would provide valuable insights for fine-tuning the hybrid approach.