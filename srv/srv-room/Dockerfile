FROM golang:1.8
MAINTAINER lanziliang <lanziliang11@gmail.com>
COPY . /go/src/chess/srv/srv-room
RUN go install chess/srv/srv-room
ENTRYPOINT ["/go/bin/srv-room"]