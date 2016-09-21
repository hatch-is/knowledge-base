FROM golang:alpine

MAINTAINER Vanya Andreychuk "vanya@tep.io"

COPY .cache /go/src

RUN apk add --update curl 

RUN go install knowledge-base

RUN rm -rf /go/src/*

CMD ["/go/bin/knowledge-base"]

HEALTHCHECK CMD curl -f 127.0.0.1:3810 || exit 1

EXPOSE 3810
