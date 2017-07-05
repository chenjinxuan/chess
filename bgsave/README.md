# bgsave
[![Build Status](https://travis-ci.org/gonet2/bgsave.svg?branch=master)](https://travis-ci.org/gonet2/bgsave)

Dump records from redis into mongodb.    
The records will be written to mongodb ASAP, callers should control the frequency of writing. 

The format of the record is defined as :          
key(tablname:record_id) -> value(packed with msgpack, optional snappy compression)    

# environment variables:
* REDIS_HOST : eg: 127.0.0.1:6379    
* MONGODB_URL : eg: mongodb://172.17.42.1/mydb
* NSQD_HOST :  eg: http://172.17.42.1:4151
* ENABLE_SNAPPY: eg: true

# install
install gpm, gvp first        
$go get -u https://github.com/gonet2/bgsave/        
$cd bgsave     
$source gvp        
$gpm       
$go install bgsave         

#install with docker
docker build -t bgsave .     
