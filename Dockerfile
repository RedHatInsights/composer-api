FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY main.go .

RUN go build -o composer-api main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/composer-api .

CMD ["./composer-api"]
