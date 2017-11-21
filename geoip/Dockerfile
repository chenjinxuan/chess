FROM golang:1.7.6
MAINTAINER chenjinxuan <jinxuanchen666@163.com>
COPY . /go/src/chess/geoip
RUN go install chess/geoip
ENTRYPOINT ["/go/bin/geoip"]
EXPOSE 50000
