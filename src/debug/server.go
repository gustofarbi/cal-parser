package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"io/ioutil"
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
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
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
	println(os.Getenv("MINIO_SECRET_KEY"))
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
