FROM golang:1.13 as builder

WORKDIR /build

COPY go.mod go.sum ./
COPY app app
COPY vendor vendor

RUN GOOS=linux GOARCH=amd64 GOFLAGS=-mod=vendor go build -o sql-alert-manager -a ./app

FROM alpine:3.10

# Update certificate
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /build/app /sql-alert-manager

CMD ["/sql-alert-manager"]
