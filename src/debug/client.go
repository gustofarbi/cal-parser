package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "svg/svg/generated"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	client := pb.NewCalendarRendererClient(conn)
	filepath := pb.Filepath{Path: "examples/wandkalender_a3-hoch_month.svg"}
	status, err := client.RenderCalendar(context.Background(), &filepath)
	if err != nil {
		println(err.Error())
	}
	if status != nil {
		fmt.Printf("status code: %d , msg: %s", status.GetCode(), status.GetMsg())
	}

	defer conn.Close()
}
