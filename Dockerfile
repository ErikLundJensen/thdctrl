# First stage: Build the Go project
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o thdctrl .

# Second stage: Run the built executable
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built executable from the builder stage
COPY --from=builder /app/thdctrl .

# Command to run the executable
CMD ["./thdctrl"]
