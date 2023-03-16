# Use the official Golang image as the base image
FROM golang:1.17-alpine as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o qbt-proxy .

# Use the alpine image for a smaller final image
FROM alpine:latest

# Install ca-certificates to enable HTTPS support
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/qbt-proxy /usr/local/bin/qbt-proxy

# Set the working directory
WORKDIR /root

# Expose the default HTTP and HTTPS ports
EXPOSE 9487 9443

# Run the qbt-proxy binary
CMD ["qbt-proxy"]

