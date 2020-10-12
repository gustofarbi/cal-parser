package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"log"
	"path"
	"svg/debug/store"
	pb "svg/svg/generated"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		println(err.Error())
	}
	svgPath := "examples/wandkalender_a3-hoch_month.svg"
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewCalendarRendererClient(conn)
	minioClient := store.MinioClient("localhost:9000")

	bucketName := "calendar-data"
	svgPathRemote := "foo/bar/"+path.Base(svgPath)
	_, err = minioClient.FPutObject(
		context.Background(),
		bucketName,
		svgPathRemote,
		svgPath,
		minio.PutObjectOptions{},
	)
	if err != nil {
		log.Fatalln(err)
	}
	filepath := pb.Filepath{Path: svgPathRemote, Bucket: bucketName}
	status, err := client.RenderCalendar(context.Background(), &filepath)
	if err != nil {
		println(err.Error())
	}
	if status != nil {
		fmt.Printf("status code: %d , msg: %s", status.GetCode(), status.GetMsg())
	}
}
