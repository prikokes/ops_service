#!/bin/bash

# Start PostgreSQL container
echo "Starting PostgreSQL container..."
docker-compose up -d postgres

# Wait for PostgreSQL to start
echo "Waiting for PostgreSQL to start..."
sleep 5

# Run application
echo "Starting application..."
go run cmd/app/main.go 