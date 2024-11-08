# Step 1: Build the Go app in a build stage
# Use the official Golang image to build the application
FROM golang:1.23-alpine3.19 AS build

# Install git, gcc, and other dependencies required for Go modules and CGO
RUN apk add --no-cache git gcc musl-dev sqlite-dev curl make

# Install Atlas database migration tool
RUN curl -sSfL https://atlasgo.sh | sh

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files to leverage Docker cache
COPY go.mod go.sum ./

# Download Go modules and dependencies
RUN go mod download

# Copy the rest of the application's source code
COPY . .

# Setup the database with automatic confirmation
RUN echo | make setup-db

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o myapp .

# Step 2: Create the final, minimal Alpine image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Install make, curl, and other necessary tools
RUN apk add --no-cache make curl

# Install Atlas database migration tool
RUN curl -sSfL https://atlasgo.sh | sh

# Install Go in the final image
RUN apk add --no-cache go

# Copy the compiled Go binary and the Makefile from the build stage
COPY --from=build /app/myapp /app/muslim_referrals.db /app/.env ./

# Copy the frontend build
COPY frontend_build ./frontend_build

# Set environment variables from the .env file
ENV SQLITE_DB_PATH=muslim_referrals.db

# Expose the port your application runs on (adjust this if necessary)
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]