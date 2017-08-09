#!/bin/sh
govendor update chess/common/consul
govendor update chess/common/define
govendor update chess/common/log
govendor update chess/common/services
go run *.go