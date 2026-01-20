# Build Stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
# COPY go.mod go.sum ./ # Root has no go.mod
COPY core/go.mod core/go.sum ./core/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Download dependencies
WORKDIR /app/core
RUN go mod download

# Copy Source
WORKDIR /app
COPY pkg ./pkg
COPY core ./core

# Build
WORKDIR /app/core
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go

# Runtime Stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /server .
COPY core/config ./config

EXPOSE 8080
CMD ["./server"]
