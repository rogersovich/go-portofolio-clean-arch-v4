# Step 1: Build the Go application
FROM golang:1.23.3 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application (main.go is in the cmd directory)
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/app cmd/main.go

# Step 2: Create a minimal container for running the application
FROM alpine:latest

# Install any necessary dependencies (e.g., certificates)
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder container
COPY --from=builder /app/bin/app .

# Set the entry point for the container (the built Go binary)
ENTRYPOINT ["./app"]

# Expose the port on which your application will run
EXPOSE 4000
