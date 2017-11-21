FROM golang:1.8
MAINTAINER chenjinxuan <jinxuanchen666@163.com>
COPY . /go/src/chess/srv/srv-chat
RUN echo "Asia/Shanghai" > /etc/timezone & dpkg-reconfigure -f noninteractive tzdata & go install chess/srv/srv-chat
ENTRYPOINT ["/go/bin/srv-chat"]
RUN mkdir /data
VOLUME /data
