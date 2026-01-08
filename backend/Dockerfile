FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy

RUN go build -o backend

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/backend /app/backend

EXPOSE 8080

ENTRYPOINT ["/app/backend"]
