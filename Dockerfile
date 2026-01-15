# --- Build Stage ---
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates # Ensure ca-certificates package is explicitly installed

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application binary
# Use CGO_ENABLED=0 for statically linked binary (better for minimal base images)
# Use -ldflags="-s -w" to strip debug info and reduce size
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o serve-player-data ./cmd/serve


# --- Final Stage ---
# Use a minimal base image like scratch or a distroless image
FROM scratch

# Copy only the ca-certificates needed for HTTPS
# (Required if your app makes outgoing HTTPS calls)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/serve-player-data /usr/local/bin

# Command to run the executable when the container starts
CMD ["/usr/local/bin/serve-player-data"]

