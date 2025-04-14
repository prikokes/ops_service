FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Install necessary tools
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

# Expose the application port
EXPOSE 8080 3000 9000

# Run the application
CMD ["./app"] 