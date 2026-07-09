FROM golang:1.26.3-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /todoapp ./cmd/todoapp

FROM alpine:3.23.5

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /todoapp /app/todoapp
COPY --from=builder app/public /app/public

ENV PROJECT_ROOT=/app

CMD ["/app/todoapp"]
