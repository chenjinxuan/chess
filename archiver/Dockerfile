FROM golang:latest
MAINTAINER xtaci <daniel820313@gmail.com>
COPY . /go/src/archiver
RUN go install archiver
RUN go install archiver/replay
RUN mkdir /data
VOLUME /data
