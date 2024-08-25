# BUild statge
FROM golang:1.23 AS builder

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download



# Copy the code into the container
COPY /*.go ./
COPY ordered_map/*.go ./ordered_map/
COPY rabbit/*.go ./rabbit/
COPY handler/*.go ./handler/



# Build the application
RUN go build -o app ./app.go

# Run stage
FROM alpine
WORKDIR /
VOLUME /OUTPUT

COPY --from=builder /build/app /app

RUN apk add bash
RUN apk add curl


# Command to run the executable
CMD ["/app"]
