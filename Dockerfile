FROM golang:1.17.8-alpine3.15 as builder

RUN apk add --update --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app

COPY go.mod .

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o comeet -ldflags="-s -w" main.go

# ---------------------------------------------

FROM scratch

COPY --from=builder /app/comeet /usr/bin/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/usr/bin/comeet"]