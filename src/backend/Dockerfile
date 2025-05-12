FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o element-api/server element-api/server.go

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/element-api/server ./element-api/server
VOLUME /app/scrap
EXPOSE 8888
WORKDIR /app/element-api
CMD ["./server"]
