FROM golang:latest
MAINTAINER xtaci <daniel820313@gmail.com>
COPY . /go/src/github.com/xtaci/chat
RUN go install github.com/xtaci/chat
ENTRYPOINT ["/go/bin/chat"]
EXPOSE 10000 8080 6060
RUN mkdir /data
VOLUME /data
