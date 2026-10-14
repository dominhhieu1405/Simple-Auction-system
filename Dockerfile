# syntax=docker/dockerfile:1

# Image chỉ đóng gói binary đã build sẵn trong dist/, không compile.
FROM alpine:3.24.2

ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup \
    && mkdir -p /app/data \
    && chown -R appuser:appgroup /app

WORKDIR /app

COPY --chown=appuser:appgroup --chmod=755 \
    dist/auction-server-${TARGETOS}-${TARGETARCH} \
    ./auction-server

ENV DB_PATH=/app/data/auction.db \
    SERVER_PORT=8080

USER appuser

EXPOSE 8080
VOLUME /app/data

HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
    CMD wget -q --spider "http://127.0.0.1:${SERVER_PORT}/api/auctions" || exit 1

ENTRYPOINT ["./auction-server"]