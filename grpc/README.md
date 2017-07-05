### 定义 gRPC service 和方法 request 以及 response 的类型


- 安装grpc `go get google.golang.org/grpc`
- 安装protobuf `brew install protobuf`
- 编写proto文件，参考[Language Guide(proto3 )](https://developers.google.com/protocol-buffers/docs/proto3#packages-and-name-resolution)
- 生成客户端和服务器端代码， eg: `./codegen.sh auth/auth.proto`