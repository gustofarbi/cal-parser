package remote

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"os"
)

type MinioDir struct {
	c      *minio.Client
	bucket string
}

func New(bucket string) *MinioDir {
	accessKeyId := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	endpoint := os.Getenv("MINIO_URL")
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
	})
	if err != nil {
		panic(err)
	}
	return &MinioDir{c: client, bucket: bucket}
}

func (m *MinioDir) ListObjects(dir string) ([]string, error) {
	objectInfoChan := m.c.ListObjects(context.Background(), m.bucket, minio.ListObjectsOptions{Prefix: dir})
	var files []string
	for objectInfo := range objectInfoChan {
		if objectInfo.Err != nil {
			panic(objectInfo.Err)
		}
		files = append(files, objectInfo.Key)
	}
	return files, nil
}

func (m *MinioDir) PutObject(path string, r io.Reader) error {
	_, err := m.c.PutObject(context.Background(), m.bucket, path, r, 0, minio.PutObjectOptions{})
	return err
}

func (m *MinioDir) GetObject(path string) (io.Reader, error) {
	return m.c.GetObject(context.Background(), m.bucket, path, minio.GetObjectOptions{})
}
