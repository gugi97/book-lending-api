FROM golang:1.23 AS builder

WORKDIR /app

# Pre copy go modules definitions
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code
COPY . .

# Download any additional modules that may have been introduced after copying source
RUN go mod download

# Build the binary with disabled CGO for static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/server /app/server

# Expose port 8080 in container
EXPOSE 8080

CMD ["/app/server"]