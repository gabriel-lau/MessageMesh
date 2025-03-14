FROM golang:1.23.7-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git build-base

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o messagemesh .

# Create a smaller final image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/messagemesh /app/

# Set environment variables
ENV HEADLESS=true

# Create a .env file with the required environment variables
RUN echo "HEADLESS=true" > /app/.env && \
    echo "USERNAME=docker" >> /app/.env

# Expose any necessary ports (adjust as needed)
# EXPOSE 8080

# Run the application
CMD ["/app/messagemesh"] 