# First stage: build the application
FROM golang:alpine AS builder

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY app/ ./app/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o discord2http ./app

# Second stage: final image
FROM alpine:latest

# Add ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/discord2http /app/discord2http

# Copy start script from scripts folder
COPY scripts/start.sh /app/start.sh
RUN chmod +x /app/start.sh

# Run the application
CMD ["/app/start.sh"]
