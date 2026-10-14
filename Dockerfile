FROM golang:1.27.2-alpine3.24 AS builder

WORKDIR /source

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -buildvcs=false \
    -trimpath \
    -ldflags="-s -w" \
    -o /output/auction-server \
    .

FROM alpine:3.24.2

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup \
    && mkdir -p /app/data \
    && chown -R appuser:appgroup /app

WORKDIR /app

COPY --from=builder \
    --chown=appuser:appgroup \
    /output/auction-server \
    ./auction-server

USER appuser

EXPOSE 8080

ENTRYPOINT ["./auction-server"]