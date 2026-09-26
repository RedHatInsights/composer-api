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

ARG MIGRATE_VERSION=v4.18.3
                                                
RUN microdnf install -y tar gzip jq && microdnf clean all && \
    curl -fsSL https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz \
    | tar -xz -C /usr/local/bin migrate

WORKDIR /app

COPY --from=builder /app/composer-api .
COPY internal/database/migrations/ /migrations/
COPY deploy/migrate.sh /migrate.sh

EXPOSE 8080

USER 1001

CMD ["./composer-api"]
