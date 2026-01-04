ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN make generate \
  make build-docker


FROM debian:bookworm
COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/static /static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["run-app"]
