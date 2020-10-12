package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"
	"svg/debug/store"
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

	start := time.Now()
	minioClient := store.MinioClient("localhost:9000")
	obj, err := minioClient.GetObject(
		context.Background(),
		filepath.Bucket,
		filepath.Path,
		minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(obj)
	if err != nil {
		log.Fatalln(err)
	}

	var foo svg.Svg
	err = xml.Unmarshal(buf.Bytes(), &foo)
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	// todo: unite src and node-mapping

	year := time.Now().Year() +1
	cal := svg.NewCalendar(buf.Bytes())
	dims := strings.Split(foo.ViewBox, " ")
	size := 2000.0
	widthViewbox, _ := strconv.ParseFloat(dims[2], 64)
	heightViewbox, _ := strconv.ParseFloat(dims[3], 64)
	scalingRatio := size / widthViewbox

	cal.Parse(foo, buf.String(), scalingRatio)
	f, err := ioutil.TempFile("/tmp", "*.svg")
	if err != nil {
		log.Fatalln(err)
	}
	prefix := strings.Replace(filepath.Path, path.Base(filepath.Path), "", -1)
	for month := 1; month <= 12; month++ {
		println("rendering month: " + strconv.Itoa(month))
		im := cal.Render(year, month, size, size*(heightViewbox/widthViewbox))
		im.SavePNG(f.Name())
		calendarPath := fmt.Sprintf("%s%d_%d.png", prefix, year, month)
		minioClient.FPutObject(
			context.Background(),
			filepath.Bucket,
			calendarPath,
			f.Name(),
			minio.PutObjectOptions{})
	}

	return &pb.Status{
		Msg:  fmt.Sprintf("done in: %vs\n", time.Since(start).Seconds()),
		Code: 0,
	}, nil
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
