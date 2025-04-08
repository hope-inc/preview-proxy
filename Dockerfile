FROM golang:1.24-bookworm AS builder

ADD . /app
WORKDIR /app

RUN make preview-proxy

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /opt/preview-proxy
COPY --from=builder /app/preview-proxy /opt/preview-proxy/app

CMD ["/opt/preview-proxy/app"]
