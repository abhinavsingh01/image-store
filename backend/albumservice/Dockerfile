FROM golang:1.21 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /albumservice

COPY scripts/wait.sh /usr/local/bin/wait.sh

CMD ["/usr/local/bin/wait.sh", "mysql:3306", "--", "/albumservice"]
