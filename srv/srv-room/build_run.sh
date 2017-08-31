#!/bin/sh
govendor update chess/common/config
govendor update chess/common/consul
govendor update chess/common/db
govendor update chess/common/define
govendor update chess/common/helper
govendor update chess/common/log
govendor update chess/common/services
govendor update chess/models

docker rm -f room-1
docker build --no-cache --rm=true -t room .
docker run --rm=true -it -p 14001:14001 -p 14101:14101 \
	--env-file ./.env \
	--name room-1 \
	room --address 192.168.40.157 --port 14001 --check-port 14101 --service-id room-1
