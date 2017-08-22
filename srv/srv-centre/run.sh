#!/bin/sh
govendor update chess/common/config
govendor update chess/common/consul
govendor update chess/common/db
govendor update chess/common/define
govendor update chess/common/helper
govendor update chess/common/log
govendor update chess/common/services
go run *.go --address 192.168.40.157 --port 10001 --service-id centre-1