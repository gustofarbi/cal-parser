package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

var (
	port = flag.Int("port", 80, "port number to listen on")
)

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(fmt.Sprintf("cannot connect on port %d: %s", *port, err))
	}

	server := grpc.NewServer()

	server.Serve(listener)
}
