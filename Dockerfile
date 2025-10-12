## Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Pre-fetch dependencies
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

# Copy source
COPY . .

# Build static binary (embeds migrations via Go embed)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /price-story ./cmd/price-story

## Runtime stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 appuser
WORKDIR /app

# Copy built binary
COPY --from=builder /price-story ./price-story

# App config
ENV PORT=8080
EXPOSE 8080
USER appuser

CMD ["./price-story"]