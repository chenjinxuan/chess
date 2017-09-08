FROM golang:1.8
MAINTAINER lanziliang <lanziliang11@gmail.com>
COPY . /go/src/chess/srv/srv-auth
RUN echo "Asia/Shanghai" > /etc/timezone & dpkg-reconfigure -f noninteractive tzdata & go install chess/srv/srv-auth
ENTRYPOINT ["/go/bin/srv-auth"]
