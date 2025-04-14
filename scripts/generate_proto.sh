#!/bin/bash

set -e

# Create directories if they don't exist
mkdir -p internal/grpc

# Generate gRPC code
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  api/proto/pvz.proto

# Move generated files to internal/grpc
mv api/proto/pvz.pb.go internal/grpc/
mv api/proto/pvz_grpc.pb.go internal/grpc/

echo "Proto files generated successfully!" 