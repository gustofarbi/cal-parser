package main

import (
	"encoding/xml"
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"svg/svg"
	"sync"
)

func main() {
	defer imagick.Terminate()
	var foo svg.Svg
	data, err := ioutil.ReadFile("wandkalender_a4-hoch_month.svg")
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	err = xml.Unmarshal(data, &foo)
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}

	c := svg.NewCalendar()
	c.Parse(foo, string(data), 1.5)

	return
	//fmt.Printf("%#v\n", foo)

	pngFile := "foo.png"
	svgFile := "wandkalender_a3-hoch_month.svg"
	cmd := exec.Command("rsvg-convert", "-w", "842", "-h", "1190", svgFile, "-o", pngFile)
	cmd.Run()

	var wg sync.WaitGroup
	os.Mkdir("images", 0775)
	var lock sync.Mutex
	for i := 0; i < 100; i++ {
		go func(i int) {
			wg.Add(1)
			defer wg.Done()
			lock.Lock()
			wand := imagick.NewMagickWand()
			defer wand.Destroy()
			pixel := imagick.NewPixelWand()
			defer pixel.Destroy()
			draw := imagick.NewDrawingWand()
			defer draw.Destroy()
			lock.Unlock()
			pixel.SetColor("white")
			draw.SetFontSize(60.0)
			err := wand.NewImage(2000, 2000, pixel)
			if err != nil {
				fmt.Println(err)
			}
			draw.Annotation(100, 100, "foobar")
			err = wand.DrawImage(draw)
			if err != nil {
				fmt.Println(err)
			}
			err = wand.WriteImage("images/" + strconv.Itoa(i) + ".jpg")
			if err != nil {
				fmt.Println(err)
			}
		}(i)
	}
	wg.Wait()
	files, _ := ioutil.ReadDir("images")
	fmt.Println("file count: " + strconv.Itoa(len(files)))
	//wand := imagick.NewMagickWand()
	//wand.SetResolution(842, 1190)
	//err = wand.ReadImage("wandkalender_a3-hoch_month.svg")
	//wand.SetImageFormat("png")
	//wand.SetFormat("png")
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err = wand.WriteImage("foo.jpg")
	//err = wand.WriteImage("foo.png")
}
