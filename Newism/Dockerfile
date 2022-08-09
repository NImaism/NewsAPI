
# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="NImaism"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]