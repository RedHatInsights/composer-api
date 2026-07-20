FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder

WORKDIR /app

COPY main.go .

RUN go build -o composer-api main.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest 

WORKDIR /app

COPY --from=builder /app/composer-api .

CMD ["./composer-api"]
