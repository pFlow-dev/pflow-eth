# Start from the official golang image
FROM golang:1.21 as builder

LABEL maintainer="stackdump <myork@stackdump.com>"

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=1 GOOS=linux go build \
-ldflags="-extldflags=-static" \
-tags sqlite_omit_load_extension \
-a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside
EXPOSE 8080

# Command to run the executable
CMD ["./main"]