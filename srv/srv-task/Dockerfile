FROM golang:1.8
MAINTAINER lanziliang <lanziliang11@gmail.com>
COPY . /go/src/chess/srv/srv-task
RUN go install chess/srv/srv-task
ENTRYPOINT ["/go/bin/srv-task"]