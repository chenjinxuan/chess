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

docker rm -f task-2
docker build --no-cache --rm=true -t task .
docker run --rm=true -it -p 15002:15002 -p 15102:15102 \
	--env-file ./.env \
	--name task-2 \
	task --address 192.168.60.164 --port 15002 --check-port 15102 --service-id task-2
