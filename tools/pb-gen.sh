#!/bin/sh
# arch       auth       bgsave     chat       game       geoip      libs       pb-gen.sh  rank       snowflake  tools      wordfilter
cd ..;cd auth
protoc  ./*.proto --go_out=plugins=grpc:src/proto
cd ..;cd bgsave
protoc  ./*.proto --go_out=plugins=grpc:src/proto
cd ..;cd chat
protoc  ./*.proto --go_out=plugins=grpc:src/proto
cd ..;cd game
protoc  ./*.proto --go_out=plugins=grpc:src/proto
cd ..;cd geoip
protoc  ./*.proto --go_out=plugins=grpc:src/proto
cd ..;cd rank
protoc  ./*.proto --go_out=plugins=grpc:src/proto
cd ..;cd snowflake
protoc  ./*.proto --go_out=plugins=grpc:src/proto
cd ..;cd wordfilter
protoc  ./*.proto --go_out=plugins=grpc:src/proto
