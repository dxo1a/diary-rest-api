FROM golang:1.24.2-alpine AS builder
RUN apk add --no-cache git gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o action-service ./server.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates sqlite-libs
WORKDIR /app
COPY --from=builder /app/action-service ./

EXPOSE 8081
ENTRYPOINT [ "./action-service" ]