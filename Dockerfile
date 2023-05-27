FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY /app/go.mod .
RUN go mod download

# Copy the source code to the working directory
COPY /app .

# Build the Go application
RUN go build -o handle-events .

# Set the entry point for the container
ENTRYPOINT ["./handle-events"]