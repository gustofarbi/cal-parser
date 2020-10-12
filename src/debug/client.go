package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"log"
	"path"
	pb "svg/svg/generated"
)

func main() {
	svgPath := "examples/wandkalender_a3-hoch_month.svg"
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewCalendarRendererClient(conn)
	minioClient := MinioClient()

	bucketName := "calendar-data"
	_, err = minioClient.FPutObject(
		context.Background(),
		bucketName,
		path.Base(svgPath),
		svgPath,
		minio.PutObjectOptions{},
	)
	if err != nil {
		log.Fatalln(err)
	}
	filepath := pb.Filepath{Path: svgPath, Bucket: bucketName}
	status, err := client.RenderCalendar(context.Background(), &filepath)
	if err != nil {
		println(err.Error())
	}
	if status != nil {
		fmt.Printf("status code: %d , msg: %s", status.GetCode(), status.GetMsg())
	}
}
