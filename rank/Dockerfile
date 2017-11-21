FROM golang:latest
MAINTAINER chenjinxuan <jinxuanchen666@163.com>
COPY . /go/src/rank
RUN go install rank
ENTRYPOINT /go/bin/rank
RUN mkdir /data
VOLUME /data
EXPOSE 50001
