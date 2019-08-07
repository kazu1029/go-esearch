FROM golang:1.12.0-alpine3.9

WORKDIR /go/src/app

COPY ./search-api .
ENV GO111MODULE=on

RUN apk add --no-cache \
      alpine-sdk \
      git \
    && go get github.com/pilu/fresh

CMD ["fresh"]
