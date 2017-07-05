#!/bin/sh
#上传numbers.xlsx到etcd
FILE=numbers.xlsx
ETCD_HOST=192.168.99.100:2379

if [ "$(uname)" == "Darwin" ]; then
	base64 $FILE > $FILE.base64
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
	base64 -w 0 $FILE > $FILE.base64
fi

curl http://$ETCD_HOST/v2/keys/numbers -XPUT --data-urlencode value@$FILE.base64
rm $FILE.base64
