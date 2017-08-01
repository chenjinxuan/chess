FROM golang:latest
MAINTAINER lanziliang <lanziliang11@gmail.com>
COPY . /go/src/room
RUN go install room
ENTRYPOINT ["/go/bin/room"]
EXPOSE 10000
