FROM golang:1.21.5-bullseye as builder
WORKDIR /build
RUN apt-get update && \
    apt-get install -y --no-install-recommends make git && \
    go env -w GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM debian:bullseye
WORKDIR /app
COPY --from=builder /build/bin .
COPY config ./config
EXPOSE 80
ENTRYPOINT ["/app/elune-backend"]