FROM golang:1.8
MAINTAINER lanziliang <lanziliang11@gmail.com>
COPY . /go/src/chess/api
RUN go install chess/api
ENTRYPOINT ["/go/bin/api"]