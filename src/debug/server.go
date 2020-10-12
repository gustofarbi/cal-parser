package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"svg/svg"
	_ "svg/svg"
	pb "svg/svg/generated"
	"time"
)

var (
	port = flag.Int("port", 50051, "port number to listen on")
)

type calendarRenderer struct {
	pb.UnimplementedCalendarRendererServer
}

func (c *calendarRenderer) RenderCalendar(ctx context.Context, filepath *pb.Filepath) (*pb.Status, error) {
	println("got request: " + filepath.Path)
	base := path.Base(filepath.Path)
	ext := path.Ext(filepath.Path)
	println(base, ext)

	start := time.Now()
	var foo svg.Svg
	data, err := ioutil.ReadFile(filepath.Path)
	minioClient := MinioClient()
	f, err := ioutil.TempFile("/tmp", "*.svg")
	if err != nil {
		log.Fatalln(err)
	}
	err = minioClient.FGetObject(
		context.Background(),
		filepath.Bucket,
		filepath.Path,
		f.Name(),
		minio.GetObjectOptions{},
	)
	if err != nil {
		log.Fatalln(err)
	}

	accessKeyId := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	client, err := minio.New("store:9000", &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatalln("cannot connect to minio: " + err.Error())
	}
	err = client.FGetObject(context.Background(),
		"foobar",
		path.Base(filepath.Path),
		filepath.Path,
		minio.GetObjectOptions{})

	//err = client.MakeBucket(context.Background(), "foobar", minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		log.Fatalf("shit happened: %s\n", err)
	}
	err = xml.Unmarshal(data, &foo)
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	// todo: unite src and node-mapping

	year, m, _ := time.Now().Date()
	month := int(m) - 8
	year += 1
	cal := svg.NewCalendar(data)
	dims := strings.Split(foo.ViewBox, " ")
	size := 2000.0
	widthViewbox, _ := strconv.ParseFloat(dims[2], 64)
	heightViewbox, _ := strconv.ParseFloat(dims[3], 64)
	scalingRatio := size / widthViewbox

	cal.Parse(foo, string(data), scalingRatio)
	cal.Render(year, month, size, size*(heightViewbox/widthViewbox))

	msg := fmt.Sprintf("done in: %vs\n", time.Since(start).Seconds())

	status := pb.Status{}
	status.Code = 0
	status.Msg = msg
	return &status, nil
}

func main() {
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		println(err.Error())
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(fmt.Sprintf("cannot connect on port %d: %s", *port, err))
	}

	server := grpc.NewServer()
	pb.RegisterCalendarRendererServer(server, &calendarRenderer{})
	go func() {
		http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintln(writer, "hello buddy")
		})
		println("http server on port: 8080")
		http.ListenAndServe(":8080", nil)
	}()

	println("grpc server on port: " + strconv.Itoa(*port))
	server.Serve(listener)
}
