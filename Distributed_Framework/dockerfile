# Use the latest stable version of the official Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files over to the container
COPY go.mod ./
COPY go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of your application code over to the container
COPY . ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/gomongodb .

# Expose port 8080 for the application
EXPOSE 8080

# Command to run the executable
CMD ["./bin/gomongodb"]
