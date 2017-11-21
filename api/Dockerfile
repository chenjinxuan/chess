FROM golang:1.8
MAINTAINER chenjinxuan <jinxuanchen666@163.com>
COPY . /go/src/chess/api
RUN echo "Asia/Shanghai" > /etc/timezone & dpkg-reconfigure -f noninteractive tzdata & go install chess/api
ENTRYPOINT ["/go/bin/api"]
