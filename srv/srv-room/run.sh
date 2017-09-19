#!/bin/sh
govendor update chess/common/config
govendor update chess/common/consul
govendor update chess/common/db
govendor update chess/common/define
govendor update chess/common/helper
govendor update chess/common/log
govendor update chess/common/services
govendor update chess/models
go run *.go --address 192.168.40.157 --port 14001 --check-port 14101 --service-id room-1 --services centre --services task --services sts --services chat