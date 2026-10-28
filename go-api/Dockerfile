# ---------- Build stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache de dependencias
COPY go.mod ./
RUN go mod download || true

COPY . .
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/go-api .

# ---------- Runtime stage ----------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget

WORKDIR /app
COPY --from=builder /out/go-api ./go-api

EXPOSE 8080
CMD ["./go-api"]