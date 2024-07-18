# Start with a base Golang Alpine image
FROM golang:alpine AS build

# Install necessary dependencies
RUN apk update && apk add --no-cache gcc libc-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Copy the config.json into the container
COPY config/config.json /app/config/config.json

# Build the Go application
RUN go build -o /bin/server cmd/algosync-service/main.go

# Start a new stage for a small image
FROM alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the previous stage
COPY --from=build /bin/server /bin/server

# Copy the config.json from the previous stage
COPY --from=build /app/config/config.json /app/config/config.json

# Install PostgreSQL client for running migrations
RUN apk update && apk add --no-cache postgresql-client

# Copy the migrations directory into the container
COPY migrations /migrations

# Expose port 4000 to the outside world
EXPOSE 4000

# Command to run the executable
CMD ["/bin/server"]
