FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder

USER 0

WORKDIR /app

RUN chown -R 1001:0 /app

USER 1001

COPY main.go .

RUN go build -o composer-api main.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

WORKDIR /app

COPY --from=builder /app/composer-api .

USER 1001

CMD ["./composer-api"]
