package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/chentaihan/grpcServer/example/pb"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

var addr = flag.String("addr", "localhost:9527", "the address to connect to")

type Ping struct {
	*grpc.ClientConn
}

func (ping *Ping) Push(in *pb.PingReq) {
	out := &pb.PingResq{}
	ping.Invoke(context.Background(), "/Ping/Push", in, out)
	data, _ := json.Marshal(out)
	fmt.Println(string(data))
}

func (ping *Ping) Push1(in *pb.PingReq) {
	out := &pb.PingResq{}
	ping.Invoke(context.Background(), "/Ping/Push1", in, out)
	data, _ := json.Marshal(out)
	fmt.Println(string(data))
}

func (ping *Ping) Push2(in *pb.PingReq) {
	out := &pb.PingResq{}
	ping.Invoke(context.Background(), "/Ping/Push2", in, out)
	data, _ := json.Marshal(out)
	fmt.Println(string(data))
}

func main() {

	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ping := &Ping{
		conn,
	}

	for i := 0; i < 100; i++ {
		in := &pb.PingReq{
			Msg: "ping " + strconv.Itoa(i),
		}
		ping.Push(in)
		ping.Push1(in)
		ping.Push2(in)
	}

}
