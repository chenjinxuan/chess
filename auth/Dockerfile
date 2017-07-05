FROM golang:latest
MAINTAINER xtaci <daniel820313@gmail.com>
COPY . /go/src/auth
RUN go install auth
ENTRYPOINT /go/bin/auth
EXPOSE 50006
