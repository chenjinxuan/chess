#!/bin/sh
govendor update chess/common/consul
govendor update chess/common/define
govendor update chess/common/log
govendor update chess/common/services
go run main.go service.go --address 192.168.40.157 --port 30001 --service-id chat-1