FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o report-server .


FROM alpine:latest

# RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/report-server .
RUN chmod +x ./report-server

EXPOSE 8081

ENTRYPOINT [ "/app/report-server" ]