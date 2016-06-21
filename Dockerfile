FROM golang:1.6.2-alpine

RUN apk update && apk upgrade && apk add git

RUN mkdir -p /go/src/github.com/byuoitav
ADD . /go/src/github.com/byuoitav/wait-for-restart-microservice

WORKDIR /go/src/github.com/byuoitav/wait-for-restart-microservice
RUN go get -d -v
RUN go install -v

CMD ["/go/bin/wait-for-restart-microservice"]

EXPOSE 8003
