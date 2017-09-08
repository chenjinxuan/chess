# 徳扑项目部署文档

### 使用第三方服务
- consul
- docker (需部署私有仓库)
- kafka

### 数据库
- redis
- mongodb
- mysql

### 编程语言
- php (pay项目)
- golang

## 部署流程

> 后续研究使用CI自动部署，下面以部署room服务为例

1.  构建docker镜像并推送到私有仓库
    
    ```bash
    $ docker build -t docker.airdroid.com/lanziliang/room:v_20170831 .
    $ docker push docker.airdroid.com/lanziliang/room:v_20170831
    ```
    
2.    服务器上拉取最新镜像
    
    ```bash
    $ docker pull docker.airdroid.com/lanziliang/room:v_20170831
    ```
    
3.  创建.env文件，写入consul配置

    ```aidl
    CONSUL_ADDRESS=59.57.13.156:8510
    CONSUL_DATACENTER=dev
    CONSUL_TOKEN=
    CONSUL_PROXY=
    ```
    
4.  使用docker启动服务

    ```bash
    $ docker run -d \
        --restart=always \
        --env-file /path/to/.env \
        -p 14001:14001 -p 14101:14101 \
        --name room-1 \
        docker.airdroid.com/lanziliang/room:v_20170831 \
        --address 59.57.13.156 \ 
        --port 14001 \
        --check-port 14101 \
        --service-id room-1 \
        --services centre \
        --services chat
    ```

 ### 各服务启动命令
 
 > 需要用到kafka服务，chat会保留一定数量的消息在内存中（默认128条），这个消息队列会定期持久化到本地磁盘，以便重启时候加载。 持久化采用boltdb，零配置, 数据存储在 VOLUME /data。 
 
```
--address       注册到consul的外部服务访问ip
--port          注册到consul的外部服务访问端口（服务监听端口）
--check-port    consul健康检查端口
--services      使用到服务，用于consul服务发现
```

    
1. agent - 网关服务

    ```bash
    $ docker run -d \
        --restart=always \
        --env-file /path/to/.env \
        -p 8898:8898 -p 8899:8899 \
        --name agent-1 \
        docker.airdroid.com/lanziliang/agent:tag \
        --tcp-listen :8898 \
        --ws-listen :8899 \
        --services room \
        --services auth
    ```
    
2. api - HTTP服务

	```bash
	$ docker run -d \
		--restart=always \
		--env-file /path/to/.env \
		-p 10086:10086 -p 10096:10096 -p 10076:10076 \
		--name api-1 \
		docker.airdroid.com/lanziliang/api:tag \
		--address 127.0.0.1 \
		--port 10086 \
		--check-port 10096  \
		--http-port 10076 \
		--service-id api-1
	```
    
3. auth - 鉴权服务

    ```bash
    $ docker run -d \
        --restart=always \
        --env-file /path/to/.env \
        -p 11001:11001 -p 11101:11101 \
        --name auth-1 \
        docker.airdroid.com/lanziliang/auth:tag \
        --address 127.0.0.1 \
        --port 11001 \
        --check-port 11101 \
        --service-id auth-1
    ```
    
4. centre - 游戏中心服， 管理房间在线人数

    ```bash
    $ docker run -d \
        --restart=always \
        --env-file /path/to/.env \
        -p 12001:12001 -p 12101:12101 \
        --name centre-1 \
        docker.airdroid.com/lanziliang/centre:tag \
        --address 127.0.0.1 \
        --port 12001 \
        --check-port 12101 \
        --service-id centre-1
    ```    
    
5. chat - 聊天服

    ```bash
    $ docker run -d \
        --restart=always \
        --env-file /path/to/.env \
        -v /data:/data \
        -p 13001:13001 -p 13101:13101 \
        --name chat-1 \
        docker.airdroid.com/lanziliang/chat:tag \
        --kafka-brokers 192.168.40.157:9092 \
        --boltdb /data/CHAT.DAT \
        --address 127.0.0.1 \
        --port 13001 \
        --check-port 13101 \
        --service-id chat-1
    ```    

6. room - 游戏服

    ```bash
    $ docker run -d \
        --restart=always \
        --env-file /path/to/.env \
        -p 14001:14001 -p 14101:14101 \
        --name room-1 \
        docker.airdroid.com/lanziliang/room:tag \
        --address 127.0.0.1 \
        --port 14001 \
        --check-port 14101 \
        --service-id room-1 \
        --services centre \
        --services task \
        --services chat
    ```    
    
7. task - 任务系统服务

    ```bash
    $ docker run -d \
        --restart=always \
        --env-file /path/to/.env \
        -p 15001:15001 -p 15101:15101 \
        --name task-1 \
        docker.airdroid.com/lanziliang/task:tag \
        --address 127.0.0.1 \
        --port 15001 \
        --check-port 15101 \
        --service-id task-1
    ```       