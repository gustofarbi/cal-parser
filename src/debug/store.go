package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func MinioClient() *minio.Client {
	accessKeyId := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	endpoint := os.Getenv("MINIO_URL")
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
	})
	handleErr(err)
	return client
}

func main() {
	err := godotenv.Load()
	handleErr(err)
	accessKeyId := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := "foobar"

	client, err := minio.New("localhost:9000", &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
	})
	handleErr(err)
	exists, err := client.BucketExists(context.Background(), bucketName)
	handleErr(err)
	if !exists {
		handleErr(client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "us-east-1"}))
	}
	uploadInfo, err := client.FPutObject(
		context.Background(),
		bucketName,
		"bar/foo.png",
		"../foo.png",
		minio.PutObjectOptions{ContentType: "application/pdf"})
	handleErr(err)
	println(uploadInfo.Size)
}

func handleErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
