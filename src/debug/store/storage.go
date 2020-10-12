package store

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func MinioClient(endpoints ...string) *minio.Client {
	accessKeyId := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	var endpoint string
	if len(endpoints) == 0 {
		endpoint = os.Getenv("MINIO_URL")
	} else {
		endpoint = endpoints[0]
	}
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
