FROM golang:1.24.2-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o metrics-service ./server.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/metrics-service ./

EXPOSE 8084
ENTRYPOINT [ "./metrics-service" ]