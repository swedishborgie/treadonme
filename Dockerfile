FROM docker.io/library/golang:1.18-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o treadonme ./webserver

FROM docker.io/library/alpine:3.15
RUN apk add --no-cache bluez
COPY --from=builder /build/treadonme /usr/bin/treadonme
EXPOSE 8089
ENV TREAD_BIND_ADDRESS=:8089
ENV TREAD_MAC_ADDRESS=""
ENV TREAD_CONNECT_TIMEOUT=60s
ENTRYPOINT [ "/usr/bin/treadonme" ]