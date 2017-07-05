#!/bin/sh
docker rm -f chat1
docker build --no-cache --rm=true -t chat .
docker run -d  --name chat1 -h chat-dev  -p 50008:50008 -e SERVICE_ID=chat1 chat
