FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/. .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .

RUN adduser -D appuser

USER appuser

CMD ["./main"]