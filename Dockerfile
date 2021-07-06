FROM golang:1.16.5-alpine3.14

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

LABEL maintainer="Zub Alexandr <zzzubalex@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .. .

RUN go build -o main .

CMD ["./main"]