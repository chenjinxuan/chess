FROM golang:latest
MAINTAINER xtaci <daniel820313@gmail.com>
COPY . /go/src/bgsave
RUN go install bgsave
ENTRYPOINT /go/bin/bgsave
EXPOSE 50004
