#!/bin/sh
govendor update chess/common/config
govendor update chess/common/consul
govendor update chess/common/db
govendor update chess/common/define
govendor update chess/common/helper
govendor update chess/common/log
govendor update chess/common/services

docker rm -f centre-1
docker build --no-cache --rm=true -t centre .
docker run --rm=true -it -p 12001:12001 -p 12101:12101 \
	--env-file ./.env \
	--name centre-1 \
	centre \
	--address 192.168.40.157 \
	--port 12001 --check-port 12101 --service-id centre-1
