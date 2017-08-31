FROM golang:1.8
MAINTAINER lanziliang <lanziliang11@gmail.com>
COPY . /go/src/chess/srv/srv-auth
RUN go install chess/srv/srv-auth
ENTRYPOINT ["/go/bin/srv-auth"]
