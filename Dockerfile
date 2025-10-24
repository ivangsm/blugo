# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gob .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    dbus \
    bluez

# Create a non-root user (optional, may need root for Bluetooth)
# RUN adduser -D -u 1000 gob

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/gob .

# Run the application
# Note: This container needs access to host D-Bus and Bluetooth
# Run with: docker run --rm -it --privileged --net=host -v /var/run/dbus:/var/run/dbus gob
CMD ["./gob"]
