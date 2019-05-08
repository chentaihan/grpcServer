package main

import (
	"fmt"
	"github.com/chentaihan/grpcServer"
	"google.golang.org/grpc"
	"net"
)

const port = 9527

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
