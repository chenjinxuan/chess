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

docker rm -f task-1
docker build --no-cache --rm=true -t task .
docker run --rm=true -it -p 15001:15001 -p 15101:15101 \
	--env-file ./.env \
	--name task-1 \
	task --address 192.168.40.157 --port 15001 --check-port 15101 --service-id task-1
