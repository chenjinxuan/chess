FROM golang:1.8
MAINTAINER lanziliang <lanziliang11@gmail.com>
COPY . /go/src/chess/srv/srv-sts
RUN go install chess/srv/srv-sts
ENTRYPOINT ["/go/bin/srv-sts"]