# chat(聊天）
[![Build Status](https://travis-ci.org/xtaci/chat.svg)](https://travis-ci.org/xtaci/chat)

## 设计理念

     EndPoint:  消息收发点（对应到一个玩家的聊天，或者一个联盟聊天） 
     PubSub: 对任意EndPoint进行发布，订阅

聊天服务器并不关心一个EndPoint对应的是一个玩家，还是一个联盟，只需要一个独立的id(snowflake-id)。       
通常在玩家注册的时候，会创建一个EndPoint，创建一个联盟的时候，也会创建一个EndPoint。     
玩家登陆后，会订阅到自己的私人EndPoint和所属联盟的EndPoint，以便接受实时聊天消息。      
CHAT会保留一定数量的消息在内存中（默认128条），这个消息队列会定期持久化到本地磁盘，以便重启时候加载。       
持久化采用boltdb，零配置, 数据存储在 VOLUME /data。        

基于PubSub的聊天服务器，优点在于可以通过多个途径**同时**访问到同一个EndPoint, 例如：      
1. 游戏内     
2. 游戏提供的离线聊天工具     
3. 提供给XMPP网关       

## 使用
参考测试用例以及chat.proto文件

## 安装
参考Dockerfile
