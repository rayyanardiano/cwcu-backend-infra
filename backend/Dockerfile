# ─── Build Stage ──────────────────────────────────────────────────────────────
FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/server ./main.go

# ─── Runtime Stage ────────────────────────────────────────────────────────────
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Jakarta

COPY --from=builder /app/bin/server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
