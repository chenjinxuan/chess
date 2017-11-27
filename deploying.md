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
    $ docker build -t $imageName:v_20170831 .
    $ docker push $imageName:v_20170831
    ```
    
2.    服务器上拉取最新镜像
    
    ```bash
    $ docker pull $imageName:v_20170831
    ```
    
3.  创建.env文件，写入consul配置

    ```aidl
    CONSUL_ADDRESS=$address:8510
    CONSUL_DATACENTER=dev
    CONSUL_TOKEN=
    CONSUL_PROXY=
    ```
    
4.  使用docker启动服务

    ```bash
    $ docker run -d \
        --net=host \
        --restart=always \
        --env-file /path/to/.env \
        -p 14001:14001 -p 14101:14101 \
        --name room01 \
        $imageName:$TAG \
        --address $address \
        --port 14001 \
        --check-port 14101 \
        --service-id room01 \
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
        --net=host \
        --restart=always \
        --env-file /path/to/.env \
        -p 8898:8898 -p 8899:8899 \
        --name agent01 \
       $imageName:$TAG \
        --tcp-listen :8898 \
        --ws-listen :8899 \
        --services room \
        --services auth
    ```
    
2. api - HTTP服务

	```bash
	$ docker run -d \
	    --net=host \
		--restart=always \
		--env-file /path/to/.env \
		-p 11008:11008 \
		--name api01 \
		$imageName:$TAG \
		--address $address \
		--http-port 11008 \
		--service-id api01
	```
    
3. auth - 鉴权服务

    ```bash
    $ docker run -d \
        --net=host \
        --restart=always \
        --env-file /path/to/.env \
        -p 11001:11001 -p 11101:11101 \
        --name auth01 \
        $imageName:$TAG \
        --address $address \
        --port 11001 \
        --check-port 11101 \
        --service-id auth01
    ```
    
4. centre - 游戏中心服， 管理房间在线人数

    ```bash
    $ docker run -d \
        --net=host \
        --restart=always \
        --env-file /path/to/.env \
        -p 12001:12001 -p 12101:12101 \
        --name centre01 \
       $imageName:$TAG \
        --address $address \
        --port 12001 \
        --check-port 12101 \
        --service-id centre01
    ```    
    
5. chat - 聊天服

    ```bash
    $ docker run -d \
        --net=host \
        --restart=always \
        --env-file /path/to/.env \
        -v /data:/data \
        -p 13001:13001 -p 13101:13101 \
        --name chat01 \
        $imageName:$TAG \
        --kafka-brokers $address:9092 \
        --boltdb /data/CHAT.DAT \
        --address $address \
        --port 13001 \
        --check-port 13101 \
        --service-id chat01
    ```    

6. room - 游戏服

    ```bash
    $ docker run -d \
        --net=host \
        --restart=always \
        --env-file /path/to/.env \
        -p 14001:14001 -p 14101:14101 \
        --name room01 \
        $imageName:$TAG \
        --address $address \
        --port 14001 \
        --check-port 14101 \
        --service-id room01 \
        --services centre \
        --services task \
        --services sts \
        --services chat
    ```    
    
7. task - 任务系统服务

    ```bash
    $ docker run -d \
        --net=host \
        --restart=always \
        --env-file /path/to/.env \
        -p 15001:15001 -p 15101:15101 \
        --name task01 \
        $imageName:$TAG \
        --address $address \
        --port 15001 \
        --check-port 15101 \
        --service-id task01
    ```       
    
    
8. sts - 统计服务

    ```bash
    $ docker run -d \
        --net=host \
        --restart=always \
        --env-file /path/to/.env \
        -p 16001:16001 -p 16101:16101 \
        --name sts01 \
        $imageName:$TAG \
        --address $address \
        --port 16001 \
        --check-port 16101 \
        --service-id sts01
    ```