#!/bin/sh
govendor update chess/common/consul
govendor update chess/common/define
govendor update chess/common/log
govendor update chess/common/services

docker rm -f agent-1
docker build --no-cache --rm=true -t agent .
docker run --rm=true -it -p 8898:8898 -p 8899:8899 \
	--env-file ./.env \
	--name agent-1 \
	agent --tcp-listen :8898 --ws-listen :8899
