FROM golang:1.8
MAINTAINER lanziliang <lanziliang11@gmail.com>
COPY . /go/src/chess/srv/srv-centre
RUN go install chess/srv/srv-centre
ENTRYPOINT ["/go/bin/srv-centre"]
