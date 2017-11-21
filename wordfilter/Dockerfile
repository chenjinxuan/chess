FROM golang:latest
MAINTAINER chenjinxuan <jinxuanchen666@163.com>
COPY . /go/src/wordfilter
RUN go install wordfilter
ENTRYPOINT ["/go/bin/wordfilter"]
EXPOSE 50002
