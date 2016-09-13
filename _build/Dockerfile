FROM golang:alpine

MAINTAINER Vanya Andreychuk "vanya@tep.io"

COPY .cache /go/src

RUN go install knowledge-base

RUN rm -rf /go/src/*

CMD ["/go/bin/knowledge-base"]

EXPOSE 3810
