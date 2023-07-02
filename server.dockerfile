
# Use an official Go runtime as a parent image
FROM golang:1.19-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any necessary dependencies
RUN go mod download

# Build the Go app
RUN go build -o myapp

# Expose port 8080 for the app to listen on
EXPOSE 8080

# Run the app when the container starts
CMD ["./myapp"]