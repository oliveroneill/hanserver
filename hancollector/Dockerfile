FROM golang:1.6

RUN apt-get update
RUN apt-get dist-upgrade -y
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/kellydunn/golang-geo
RUN go get github.com/gedex/go-instagram/instagram

ADD . /go/src/github.com/oliveroneill/hanserver/
WORKDIR /go/src/github.com/oliveroneill/hanserver/hancollector

RUN go install -v