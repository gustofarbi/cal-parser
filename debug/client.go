package main

import (
	"context"
	"google.golang.org/grpc"
	pb "svg/svg/generated"
)

func main() {
	conn, err := grpc.Dial("")
	if err != nil {
		panic(err)
	}
	client := pb.NewCalendarRendererClient(conn)
	filepath := pb.Filepath{Path: }
	client.RenderCalendar(context.Background(), )

	defer conn.Close()
}
