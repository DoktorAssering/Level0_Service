# Stage 1: Build the Go application
FROM golang:1.19 as builder

WORKDIR /app
COPY . .

WORKDIR /app/ServiceApp
RUN go mod download
RUN go build -o server .

# Stage 2: Setup the runtime container
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl postgresql-client

WORKDIR /root/
COPY --from=builder /app/ServiceApp/server .

EXPOSE 8080
CMD ["./server"]