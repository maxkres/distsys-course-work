FROM golang:alpine AS builder

WORKDIR /usr/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build

ENTRYPOINT ["./server"]
