FROM golang:1.8
MAINTAINER chenjinxuan <jinxuanchen666@163.com>
COPY . /go/src/chess/agent
RUN go install chess/agent
ENTRYPOINT ["/go/bin/agent"]
