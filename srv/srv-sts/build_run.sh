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

docker rm -f sts-1
docker build --no-cache --rm=true -t sts .
docker run --rm=true -it -p 16001:16001 -p 16101:16101 \
	--env-file ./.env \
	--name sts-1 \
	sts --address 192.168.60.164 --port 16001 --check-port 16101 --service-id sts-1
