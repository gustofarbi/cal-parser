package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gustofarbi/parser/store"
	"github.com/gustofarbi/parser/store/remote"
	"github.com/gustofarbi/parser/store/temp"
	"github.com/gustofarbi/parser/svg"
	"github.com/thoas/go-funk"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type renderRequest struct {
	Directory string `json:"directory"`
	Hash      string `json:"hash"`
	Year      int    `json:"year"`
	Width     int    `json:"width"`
}

const (
	bucket = "calendar-data"
)

var (
	minioClient store.Storer
)

func init() {
	minioClient = remote.New(bucket)
}

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/render", render)
	for {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		log.Println(err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{Status: "healthy"})
	if err != nil {
		panic(err)
	}
}

func render(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	body := r.Body
	defer func() {
		err := body.Close()
		if err != nil {
			log.Println("error closing body: " + err.Error())
		}
	}()
	var rr renderRequest
	err := json.NewDecoder(body).Decode(&rr)
	if err != nil {
		log.Fatalln("error decoding body: " + err.Error())
	}

	files, err := minioClient.ListObjects(rr.Directory)
	if err != nil {
		panic(err)
	}
	tempDir := temp.New()
	var svgs []string
	for _, file := range files {
		ext := path.Ext(file)
		if funk.Contains([]string{"svg", "jpg", "jpeg", "png"}, ext) {
			if ext == "svg" {
				svgs = append(svgs, file)
			}
			r, err := minioClient.GetObject(file)
			if err != nil {
				panic(err)
			}
			err = tempDir.PutObject(file, r) // todo maybe we need only the base
			if err != nil {
				panic(err)
			}
		}
	}

	for _, svgFile := range svgs {
		r, err := tempDir.GetObject(svgFile)
		if err != nil {
			panic(err)
		}
		var svgStruct svg.Svg
		err = xml.NewDecoder(r).Decode(&svgStruct)
		if err != nil {
			panic(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			panic(err)
		}
		year := time.Now().Year() + 1
		cal := svg.NewCalendar(b)
		dims := strings.Split(svgStruct.ViewBox, " ")
		size := 2000.0
		widthViewbox, _ := strconv.ParseFloat(dims[2], 64)
		heightViewbox, _ := strconv.ParseFloat(dims[3], 64)
		scalingRatio := size / widthViewbox

		cal.Parse(svgStruct, string(b), scalingRatio)
		prefix := rr.Hash
		for month := 1; month <= 12; month++ {
			println("rendering month: " + strconv.Itoa(month))
			im := cal.Render(year, month, size, size*(heightViewbox/widthViewbox))
			var buf bytes.Buffer
			err = png.Encode(&buf, im.Image())
			if err != nil {
				log.Println("cannot encode image: " + err.Error())
				continue
			}
			// todo locale + refinement
			calendarPath := fmt.Sprintf("%s%d_%d.png", prefix, year, month)
			err = minioClient.PutObject(calendarPath, &buf)
			if err != nil {
				log.Println("cannot write image: " + err.Error())
			}
		}
	}
	_, err = fmt.Fprintf(w, "done in %s", time.Since(start).String())
	if err != nil {
		panic(err)
	}
	log.Printf("done in %s", time.Since(start).String())
}
