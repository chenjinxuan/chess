#!/bin/sh
govendor update chess/common/consul
govendor update chess/common/define
govendor update chess/common/log
govendor update chess/common/services

docker rm -f chat-1
docker build --no-cache --rm=true -t chat .
docker run --rm=true -it -p 13001:13001 -p 13101:13101 \
	--env-file ./.env \
	-v /tmp:/data \
	--name chat-1 \
	chat \
	--kafka-brokers 192.168.40.157:9092 \
	--boltdb /data/CHAT.DAT \
	--address 192.168.40.157 \
	--port 13001 --check-port 13101 --service-id chat-1
