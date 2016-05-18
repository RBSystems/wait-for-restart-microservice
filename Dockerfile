FROM golang:1.6

RUN mkdir -p /go/src/github.com/byuoitav
ADD . /go/src/github.com/byuoitav/wait-for-reboot-microservice

WORKDIR /go/src/github.com/byuoitav/wait-for-reboot-microservice
RUN go get -d -v
RUN go install -v

CMD ["/go/bin/wait-for-reboot-microservice"]

EXPOSE 8003
