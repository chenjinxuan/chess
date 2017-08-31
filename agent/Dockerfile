FROM golang:1.8
MAINTAINER lanziliang <lanziliang11@gmail.com>
COPY . /go/src/chess/agent
RUN go install chess/agent
ENTRYPOINT ["/go/bin/agent"]
