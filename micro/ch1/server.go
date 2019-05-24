package main

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
	pd "learn_micro/micro/ch1/proto"
	"net"
)

type Server struct {
}

func (Server) SayHello(ctx context.Context, req *pd.HelloRequest) (*pd.HelloResp, error) {
	fmt.Println(req)
	return &pd.HelloResp{}, nil
}

func main() {
	req := &pd.HelloRequest{
		Name: "long",
		Contents: []string{"hello"},
	}

	bytes, err := proto.Marshal(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	// 注册grpc
	es := &Server{}
	server := grpc.NewServer()
	pd.RegisterHelloServiceServer(server, es)

	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
