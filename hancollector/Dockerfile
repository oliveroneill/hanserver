FROM golang:1.6

RUN apt-get update
RUN apt-get dist-upgrade -y

ADD . /go/src/github.com/oliveroneill/hanserver/
WORKDIR /go/src/github.com/oliveroneill/hanserver/hancollector

RUN go get ./...
RUN go install