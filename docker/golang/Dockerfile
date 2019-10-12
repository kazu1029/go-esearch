FROM golang:1.12.0-alpine3.9

WORKDIR /go/src/go-elastic/search-api

COPY . /go/src/go-elastic

ENV GO111MODULE=on

RUN apk add --no-cache \
      alpine-sdk \
      git \
    && go get github.com/pilu/fresh

CMD ["fresh"]
