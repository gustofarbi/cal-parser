package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	pb "svg/svg/generated"
)

var (
	port = flag.Int("port", 80, "port number to listen on")
)

type calendarRenderer struct {
	pb.UnimplementedCalendarRendererServer
}

func (c *calendarRenderer) RenderCalendar(ctx context.Context, filepath *pb.Filepath) (*pb.Status, error) {
	status := pb.Status{}
	status.Code = 0
	status.Msg = "rendered successfully"
	return &status, nil
}

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(fmt.Sprintf("cannot connect on port %d: %s", *port, err))
	}

	server := grpc.NewServer()
	pb.RegisterCalendarRendererServer(server, &calendarRenderer{})

	server.Serve(listener)
}
