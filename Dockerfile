FROM golang:1.24 AS builder

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/server/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /server ./server

# Expose the service's port. The value here should match the
# configured SERVER_PORT in docker-compose.yml or environment.
EXPOSE 8080

CMD ["./server"]