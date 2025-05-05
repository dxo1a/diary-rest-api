FROM golang:1.24.2-alpine AS builder
RUN apk add --no-cache git gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o api-gateway ./server.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates sqlite-libs
WORKDIR /app
COPY --from=builder /app/api-gateway ./

EXPOSE 80
ENTRYPOINT [ "./api-gateway" ]