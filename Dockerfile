FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apt-get update -qq && \
  apt-get install -no-install-recommends -y build-essential pkg-config

RUN curl -fsSL https://deb.nodesource.com/setup_current.x | bash - && build-essential

COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN make build

FROM scratch
WORKDIR /app
COPY --from=builder /app/bin .
ARG ADDR
EXPOSE 8000
CMD ["/app/api"]
