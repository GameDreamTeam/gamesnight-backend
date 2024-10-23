FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN go fmt ./...

RUN go build -o ./bin/charades ./cmd/charades

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bin/charades ./charades

CMD ["./charades"]
