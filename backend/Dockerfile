FROM golang:1.21-alpine AS builder
WORKDIR /app

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . ./

# Build backend binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /subfinder-service ./cmd/server

# Install subfinder binary compatible with Go 1.21
RUN go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@v2.6.3

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# Copy binaries
COPY --from=builder /subfinder-service /usr/local/bin/subfinder-service
COPY --from=builder /go/bin/subfinder /usr/local/bin/subfinder

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

ENV PORT=8080
ENV WORKER_COUNT=5

EXPOSE 8080
ENTRYPOINT ["subfinder-service"]
