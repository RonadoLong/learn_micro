package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"learn_micro/micro/ch3/handler"
	"learn_micro/micro/ch3/subscriber"

	example "learn_micro/micro/ch3/proto/example"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.ch3"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	_ = example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Register Struct as Subscriber
	_ = micro.RegisterSubscriber("go.micro.srv.ch3", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	_ = micro.RegisterSubscriber("go.micro.srv.ch3", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
