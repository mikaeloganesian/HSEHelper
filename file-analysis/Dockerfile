FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o file-analysis .

FROM debian:bullseye-slim

COPY --from=builder /app/file-analysis /app/file-analysis

WORKDIR /app

EXPOSE 8081

ENTRYPOINT ["./file-analysis"]
