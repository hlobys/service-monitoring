# Step 1: Build the Go application
FROM golang:1.16-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go app
RUN go build -o system-monitor

# Step 2: Create a small image to run the Go app
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/system-monitor .

# Run the Go app
CMD ["./system-monitor"]
