# Use official Golang image
FROM golang:1.23

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app inside container
RUN go build -o app .

# Expose the app port
EXPOSE 8002

# Run the app
CMD ["./app"]
