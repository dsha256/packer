#!/bin/bash
set -e

echo "Starting combined service..."

echo "Starting backend on port 3000..."
/usr/local/bin/packer &
BACKEND_PID=$!

echo "Waiting for backend to be ready..."
for i in {1..30}; do
  if curl -f http://localhost:3000/api/v1/packet/size > /dev/null 2>&1; then
    echo "Backend is ready!"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "Backend failed to start within 30 seconds"
    kill $BACKEND_PID 2>/dev/null || true
    exit 1
  fi
  echo "Waiting for backend... ($i/30)"
  sleep 1
done

echo "Starting frontend on port 80..."
cd /app/ui
exec pnpm start
