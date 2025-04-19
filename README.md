# Packet Manager

A system for managing and optimizing packet sizes for shipping.

## Docker Compose Setup

This project includes a Docker Compose configuration to run both the backend and frontend services together.

### Prerequisites

- Docker and Docker Compose installed on your system
- Git (to clone the repository)

### Running the Application

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. Access the application:
   - Frontend: http://localhost:3001
   - Backend API: http://localhost:3000

### Stopping the Application

To stop the services:
```bash
docker-compose down
```

### Rebuilding the Services

If you make changes to the code, you'll need to rebuild the services:
```bash
docker-compose up -d --build
```

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