FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN apk add --no-cache make~=4.4 && make build

FROM scratch
WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin .
COPY --from=builder /app/web/static ./web/static
COPY --from=builder /app/web/templates ./web/templates

EXPOSE $PORT_ARG
CMD ["/app/api"]
