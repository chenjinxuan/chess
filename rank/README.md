# Ranking Service(排名)
[![Build Status](https://travis-ci.org/xtaci/rank.svg)](https://travis-ci.org/xtaci/rank)

## 设计理念
对int32类型的id, score进行排名， 并用boltdb实现持久化。      
排名依据score进行，可以获得范围，比如［1，100］名的列表，可以定位某个玩家的排名，比如id为1234的排名。      
排名包含无限个集合，根据id(snowflake-id)区分，用户根据业务需求创建。          
持久化采用boltdb，零配置, 数据存储在 volume /data 。      

## 性能
采用混合策略做排名:           
1. 数据量小于1024(暂定)的时候，采用sortedset实现, 大部分操作时间复杂度为O(n)       
2. 超过之后采用rbtree实现, 时间复杂度O(logN)。         

sortedset的紧凑存储结构能充分利用cpu cache，而对rbtree的访问基本是全部cache miss;        所以必须在达到一定数据量之后，算法时间复杂度提升才能弥补cache miss.         

## 使用
参考测试用例以及rank.proto文件

## 安装
参考Dockerfile
