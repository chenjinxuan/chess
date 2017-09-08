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

docker rm -f api-1
docker build --no-cache --rm=true -t api .
docker run --rm=true -it -p 10086:10086 -p 10096:10096 -p 10076:10076 \
	--env-file ./.env \
	--name api-1 \
	api -address 192.168.60.164 --port 10086 --check-port 10096  --http-port 10076 --service-id api-1
