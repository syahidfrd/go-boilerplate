FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create app directory
WORKDIR /app

# Copy all other source code to work directory
COPY . .

# Download all the dependencies that are required
RUN go mod tidy

# Build the application
RUN go build -o binary cmd/api/main.go

# The lightweight scratch image we'll run our application within
FROM alpine:latest

# We have to copy the output from our builder stage
COPY --from=builder /app .

# Expose port
EXPOSE 8080

# Executable
ENTRYPOINT ["./binary"]
