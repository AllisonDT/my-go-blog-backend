# Stage 1: Build the Go binary.
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go binary with optimizations and static linking.
RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags="-s -w" -o app .

# Stage 2: Create a minimal image to run the binary.
FROM alpine:latest

# Install CA certificates.
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage.
COPY --from=builder /app/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
