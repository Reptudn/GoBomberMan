FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o report-server .


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/report-server .

EXPOSE 8080

ENTRYPOINT [ "/app/report-server" ]