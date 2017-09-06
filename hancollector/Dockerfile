FROM golang:1.8

RUN apt-get update
RUN apt-get dist-upgrade -y

RUN go get github.com/oliveroneill/flickgo
RUN go get github.com/gedex/go-instagram/instagram
RUN go get github.com/dghubble/oauth1
RUN go get github.com/dghubble/go-twitter/twitter
RUN go get gopkg.in/alecthomas/kingpin.v2
RUN go get github.com/nlopes/slack
RUN go get github.com/kellydunn/golang-geo
RUN go get gopkg.in/mgo.v2

ADD . /go/src/github.com/oliveroneill/hanserver/
WORKDIR /go/src/github.com/oliveroneill/hanserver/hancollector

RUN go install