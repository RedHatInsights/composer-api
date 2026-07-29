FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder

USER 0

WORKDIR /app

RUN chown -R 1001:0 /app

USER 1001

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/

RUN go build -o composer-api ./cmd/composer-api

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

WORKDIR /app

COPY --from=builder /app/composer-api .

EXPOSE 8080

USER 1001

CMD ["./composer-api"]
