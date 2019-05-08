# grpcServer

## grpc管理器

### 目的

`简化代码：原生grpc方法，每个方法都要写一堆的代码。使用grpcServer，
你只需要关注grpc方法的实现，其他的代码都不需要写`

### 原理
`通过反射实现自动注册grpc方法`

### grpc server实例
只需要写具体实现，其他的代码都不需要写
```import (
 	"context"
 	"github.com/chentaihan/grpcServer/example/pb"
 )
 
 type Ping struct {
 }
 
 //接口：/Ping/Push
 func (ping *Ping) Push(ctx context.Context, req *pb.PingReq) (*pb.PingResq, error) {
 	return &pb.PingResq{Msg: "pong:" + req.Msg}, nil
 }
 ```
 
 ### 使用实例
 ```const port = 9527
 
 func main() {
 	server := grpc.NewServer()
 	log := new(grpcServer.LogDefault)
 	rpcServer := grpcServer.NewRpcRegister(server, log)
 
 	ht := &Ping{}
 	rpcServer.Register(ht)
 
 	if err := server.Serve(createListener(port)); err != nil {
 		fmt.Println("grpc server start err", err.Error())
 	}
 }
 
 func createListener(port int) net.Listener {
 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
 	if err != nil {
 		panic(err.Error())
 	}
 
 	return lis
 }

 ```
 