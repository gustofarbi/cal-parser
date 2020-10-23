package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"path"
	"strconv"
	"svg/debug/store"
	pb "svg/svg/generated"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		println(err.Error())
	}
	//svgPath := "examples/wandkalender_a4-hoch_month.svg"
	svgPath := "examples/Architektur/tischkalender_145-145_month.svg"
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewCalendarRendererClient(conn)
	minioClient := store.MinioClient("localhost:9000")

	bucketName := "calendar-data"
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Println(err)
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln(err)
		}
	}
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := int(gen.Int63n(1000000))
	svgPathRemote := strconv.Itoa(rnd) + "/" + path.Base(svgPath)
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
