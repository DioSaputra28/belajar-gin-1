FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /

# Create bin directory
RUN mkdir -p /bin

# Copy the binary from builder to /bin/app
COPY --from=builder /app/main /bin/app

# Expose port (dynamic, will use PORT env var)
EXPOSE 8080

# Run the application
CMD ["/bin/app"]
