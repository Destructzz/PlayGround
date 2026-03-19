FROM golang:1.24-alpine AS builder
WORKDIR /app

# Pre-download modules
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source
COPY backend/. .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

FROM alpine:3.20
RUN adduser -D -g '' appuser && apk add --no-cache ca-certificates
USER appuser
WORKDIR /home/appuser
COPY --from=builder /out/api /usr/local/bin/api
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/static ./static

ENV APP_HTTP_ADDR=:8080
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/api"]
