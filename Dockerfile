# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./
EXPOSE 8000
RUN go build -o ms-starter cmd/rest/main.go

CMD ["./ms-starter"]