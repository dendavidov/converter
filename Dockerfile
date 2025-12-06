FROM golang:1.22-bookworm AS builder

WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/converter .

FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive \
    QTWEBENGINE_DISABLE_SANDBOX=1 \
    QTWEBENGINE_CHROMIUM_FLAGS="--disable-gpu --no-sandbox"

RUN apt-get update && \
    apt-get install -y --no-install-recommends calibre dbus-x11 fontconfig && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /out/converter /app/converter

ENTRYPOINT ["/app/converter"]
