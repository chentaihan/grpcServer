package main

import (
	"context"
	"github.com/chentaihan/grpcServer/example/pb"
)

type Ping struct {
}

//接口：/Ping/Push
func (ping *Ping) Push(ctx context.Context, req *pb.PingReq) (*pb.PingResq, error) {
	return &pb.PingResq{Msg: "pong:" + req.Msg}, nil
}
