FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy all source code
COPY . .

# Initialize the module and get dependencies
RUN go mod init github.com/user/subfinder-service
RUN go get github.com/gin-gonic/gin
RUN go get github.com/google/uuid
RUN go get github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o subfinder-service ./cmd/server

# Install subfinder
RUN go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest

FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/subfinder-service /usr/local/bin/subfinder-service

# Copy subfinder binary
COPY --from=builder /go/bin/subfinder /usr/local/bin/subfinder

# Create a non-root user to run the application
RUN adduser -D -g '' appuser
USER appuser

# Set environment variables
ENV PORT=8080
ENV WORKER_COUNT=5
ENV BUILD_TIMESTAMP="2025-04-03T11:02:30Z"

# Expose the port
EXPOSE 8080

# Run the application
ENTRYPOINT ["subfinder-service"]
