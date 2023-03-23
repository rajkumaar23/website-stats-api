FROM golang:1.20.2-alpine3.17

RUN apk add --no-cache git make musl-dev

RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN go build -o build

CMD "./build"
