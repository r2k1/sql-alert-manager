FROM golang:1.13 as builder

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY app app

RUN GOOS=linux GOARCH=amd64 go build -o sql-alert-manager -a ./app

FROM alpine:3.10

# Update certificates
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /build/app /sql-alert-manager

CMD ["/sql-alert-manager"]
