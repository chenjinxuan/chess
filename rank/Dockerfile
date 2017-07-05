FROM golang:latest
MAINTAINER xtaci <daniel820313@gmail.com>
COPY . /go/src/rank
RUN go install rank
ENTRYPOINT /go/bin/rank
RUN mkdir /data
VOLUME /data
EXPOSE 50001
