# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.21 as builder

# Copy local code to the container image.
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o orders

# Use a minimal image for production.
FROM alpine:3.14
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/orders-service /orders-service

# Run the service on container startup.
CMD ["/orders-service"]
