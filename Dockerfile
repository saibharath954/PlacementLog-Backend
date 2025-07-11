# Builder image
FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# Final image
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]
