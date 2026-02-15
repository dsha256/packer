# Packer

A system for managing and optimizing packet sizes for shipping.

---

# ⚠️ AI Usage

AI is used only for UI and documentation.

---

## Run Locally - Docker Compose Setup

This project includes a Docker Compose configuration to run both the backend and frontend services together.

### Prerequisites

- Docker and Docker Compose installed on your system
- Git (to clone the repository)

### Running the Application

1. Clone the repository (SSH):
   ```bash
   git clone git@github.com:dsha256/packer.git
   ```
2. Navigate to the repository root folder:
   ```bash
   cd packer
   ```

3. Start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```
   Or with Taskfile:
   ```bash
   task compose_up
   ```

4. Access the application:
   - Frontend: http://localhost:3001
   - Backend API: http://localhost:3000

### Stopping the Application

To stop the services:
```bash
docker-compose down --remove-orphans
```

### Rebuilding the Services

If you make changes to the code, you'll need to rebuild the services:
```bash
docker-compose up -d --build
```

---

## Development

### Backend (Go)

The backend is a Go application that provides API endpoints for packet management.

### Frontend (Next.js)

The frontend is a Next.js application that provides a user interface for interacting with the backend.

## Configuration

The backend configuration is stored in `config.yaml`. This file is mounted as a volume in the Docker container.

## Troubleshooting

If you encounter port conflicts, make sure no other services are using ports 3000 and 3001.

To check the logs:
```bash
docker-compose logs -f
```

To check the logs for a specific service:
```bash
docker-compose logs -f backend
docker-compose logs -f frontend
```

---

## Taskfile Commands

This project uses Taskfile for task automation. Here are the available commands:

### Development Tasks

- `task test` - Run all tests with the race flag enabled
- `task lint` - Run the Go linter to check code quality
- `task format` - Format all Go code using gofumpt and fieldalignment
- `task benchmark_packer` - Run the packer service benchmarks

### Docker Compose Tasks

- `task compose_up` - Start the Docker Compose services
- `task compose_down` - Stop and remove Docker Compose services
- `task compose_fresh_restart` - Stop all services and restart them with a fresh build

### Profiling Tasks

The following tasks launch web UIs for different types of profiling (available at http://localhost:9090):

- `task pprof_allocs_web` - Memory allocation profiling
- `task pprof_heap_web` - Heap memory profiling
- `task pprof_goroutine_web` - Goroutine profiling
- `task pprof_block_web` - Blocking profiling
- `task pprof_threadcreate_web` - Thread creation profiling
- `task pprof_trace_web` - Execution tracing
- `task pprof_profile_web` - CPU profiling
- `task pprof_symbol_web` - Symbol lookup

Note: Profiling endpoints are only available when the Go application is running. 

---

## Algorithms used for packing

### Summary of Complexities

| **Algorithm**                     | **Time Complexity**                                   | **Space Complexity**       | **Approach**                           |
|------------------------------------|------------------------------------------------------|-----------------------------|----------------------------------------|
| **`CalculateOptimalPacketsForItemsV1`** | `O((items + maxPacketSize) * len(packetSizes))`      | `O(items + maxPacketSize)` | Dynamic Programming (Backtracking)    |
| **`CalculateOptimalPacketsForItemsV2`** | `O((items + maxPacketSize) * len(packetSizes) * log(items + maxPacketSize))` | `O(items + maxPacketSize)` | Dijkstra's Algorithm with Min-Heap    |

### Key Differences in Approach

| **Factor**               | **V1 (DP)**                                | **V2 (Min-Heap / Dijkstra)**             |
|--------------------------|--------------------------------------------|------------------------------------------|
| **Algorithm Type**       | Dynamic Programming                        | Priority Queue + Greedy Traversal (Dijkstra) |
| **Main Data Structure**  | Arrays (`dpPacks`, `prevPacket`)            | Heap (`MinHeap`) + Maps (`minNumPacks`)  |
| **Backtracking**         | Uses `prevPacket` to reconstruct solution. | Uses `predecessor` map to reconstruct solution. |
| **Efficiency**           | Processes all totals up to `maxSum`.       | Prioritizes smaller totals with fewer packets first. |

### Benchmarks

- Benchmarks and other implementation details are given in this folder [internal/packer](https://github.com/dsha256/packer/tree/main/internal/packer)
- Actual benchmarks' results can be found in this folder [benchmark](https://github.com/dsha256/packer/tree/main/benchmark)