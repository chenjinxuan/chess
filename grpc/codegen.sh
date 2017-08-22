#!/usr/bin/env bash

# This script serves as an example to demonstrate how to generate the gRPC-Go
# interface and the related messages from .proto file.
#
# It assumes the installation of i) Google proto buffer compiler at
# https://github.com/google/protobuf (after v2.6.1) and ii) the Go codegen
# plugin at https://github.com/golang/protobuf (after 2015-02-20). If you have
# not, please install them first.
#
# We recommend running this script at $GOPATH/src.
#
# If this is not what you need, feel free to make your own scripts. Again, this
# script is for demonstration purpose.
#
protoc --go_out=plugins=grpc:. *.proto

# agent
cp agent.pb.go $GOPATH/src/chess/agent/proto/
cp room.pb.go $GOPATH/src/chess/agent/proto/
cp auth.pb.go $GOPATH/src/chess/agent/proto/

# srv-auth
cp agent.pb.go $GOPATH/src/chess/srv/srv-auth/proto/
cp auth.pb.go $GOPATH/src/chess/srv/srv-auth/proto/

# srv-room
cp agent.pb.go $GOPATH/src/chess/srv/srv-room/proto/
cp room.pb.go $GOPATH/src/chess/srv/srv-room/proto/
cp centre.pb.go $GOPATH/src/chess/srv/srv-room/proto/
cp chat.pb.go $GOPATH/src/chess/srv/srv-room/proto/

# srv-chat
cp chat.pb.go $GOPATH/src/chess/srv/srv-chat/proto/

# srv-centre
cp centre.pb.go $GOPATH/src/chess/srv/srv-centre/proto/

# api
cp auth.pb.go $GOPATH/src/chess/api/proto/