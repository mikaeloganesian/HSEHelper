FROM golang:1.24 AS builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o gateway .

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/gateway .

COPY --from=builder /app/docs ./docs

COPY .env .env

EXPOSE 8080

CMD ["./gateway"]
