FROM golang:1.8
MAINTAINER lanziliang <lanziliang11@gmail.com>
COPY . /go/src/chess/srv/srv-chat
RUN go install chess/srv/srv-chat
ENTRYPOINT ["/go/bin/srv-chat"]
RUN mkdir /data
VOLUME /data
