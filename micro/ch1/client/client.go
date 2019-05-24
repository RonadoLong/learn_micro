package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pd "learn_micro/micro/ch1/proto"
)

func main() {
	conn, err := grpc.Dial(":9999", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pd.NewHelloServiceClient(conn)
	resp, err := client.SayHello(context.Background(), &pd.HelloRequest{Name: "long", Contents: []string{"12323"}})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

}
