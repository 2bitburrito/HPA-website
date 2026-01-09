ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN make generate
RUN make build-docker


FROM debian:bookworm
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/static /static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["run-app"]
