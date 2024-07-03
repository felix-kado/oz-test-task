# Stage 1: Build the Go application
FROM golang:1.22-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o ozon-test ./cmd/bin

# Stage 2: Run the Go application
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ozon-test /app/ozon-test

# Add a health check to ensure the container is healthy
HEALTHCHECK CMD curl --fail http://localhost:8080/ || exit 1

# Command to run the executable
CMD ["/app/ozon-test"]
