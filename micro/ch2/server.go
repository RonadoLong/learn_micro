package main

import (
	"net"
	"net/rpc"
)

func main() {
	listener, e := net.Listen("tcp", ":9999")
	if e != nil {
		panic(e)
	}

	server := rpc.DefaultServer
	server.Accept(listener)
}
