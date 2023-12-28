FROM golang:1.21.5-alpine3.18 as builder
WORKDIR /build
RUN apk add make git && \
    go env -w GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /build/bin .
COPY config ./config
EXPOSE 80
ENTRYPOINT ["/app/elune-backend"]