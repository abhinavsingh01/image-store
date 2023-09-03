FROM golang:1.21 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /apigwservice

COPY run.sh /app/run.sh

ENTRYPOINT ["/app/run.sh"]
