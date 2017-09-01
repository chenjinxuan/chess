#!/bin/sh
govendor update chess/common/auth
govendor update chess/common/cache
govendor update chess/common/config
govendor update chess/common/consul
govendor update chess/common/db
govendor update chess/common/define
govendor update chess/common/helper
govendor update chess/common/log
govendor update chess/common/services
govendor update chess/common/storage
govendor update chess/models

docker rm -f auth-1
docker build --no-cache --rm=true -t auth .
docker run --rm=true -it -p 11001:11001 -p 11101:11101 \
	--env-file ./.env \
	--name auth-1 \
	auth \
	--address 192.168.60.164 \
	--port 11001 --check-port 11101 --service-id auth-1
