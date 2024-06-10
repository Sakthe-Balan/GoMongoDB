# Start from the official Golang image for building the application
FROM golang:1.21.4 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Set the Go environment variables
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=off \
    CGO_ENABLED=0

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/main

# Expose port 6942 to the outside world
EXPOSE 6942

# Command to run the executable
ENTRYPOINT ["/app/main"]
