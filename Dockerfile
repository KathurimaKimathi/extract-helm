# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.23.4-bullseye AS builder

# Create and change to the app directory.
WORKDIR /app

# Copy go.sum/go.mod and warm up the module cache.
COPY go.* ./

RUN go mod download

# Set the environment variable for Gin in release mode.
ENV GIN_MODE=release

# Now copy the rest of the application's source code
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server github.com/ingenium-connect/digitaltaxi

FROM debian:bullseye-slim AS production

# Install necessary packages including ca-certificates and openssl.
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    openssl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy the Go binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Set the working directory to where your binary and templates are.
WORKDIR /app

# Run the web service on container startup.
CMD ["/server"]
