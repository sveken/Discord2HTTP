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

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S discord2http -G appgroup

WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/discord2http /app/discord2http

# Copy start script from scripts folder
COPY scripts/start.sh /app/start.sh
RUN chmod +x /app/start.sh

# Change ownership of the application files
RUN chown -R discord2http:appgroup /app

# Switch to non-root user
USER discord2http

# Add OCI labels
LABEL org.opencontainers.image.title="Discord2HTTP"
LABEL org.opencontainers.image.description="This project is for turning a Discord channel, events or both into simple HTTP strings that can be parsed and formatted in games."
LABEL org.opencontainers.image.authors="Sveken"
LABEL org.opencontainers.image.url="https://github.com/sveken/Discord2HTTP"
LABEL org.opencontainers.image.source="https://github.com/sveken/Discord2HTTP"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.base.name="alpine:latest"

# Run the application
CMD ["/app/start.sh"]
