# =========== BUILDER ===========
FROM golang:1.24 AS builder

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
  && update-ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o /app/server

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates

COPY --from=builder /app/server /server

EXPOSE 8080
ENTRYPOINT ["/server"]
