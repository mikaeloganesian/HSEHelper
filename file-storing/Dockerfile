FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o file-storing .

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/file-storing .

EXPOSE 8081

CMD ["./file-storing"]
