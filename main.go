package main

import (
	"encoding/xml"
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io/ioutil"
	"svg/svg"
)

func main() {
	var foo svg.Svg
	data, err := ioutil.ReadFile("wandkalender_a3-hoch_month.svg")
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	err = xml.Unmarshal(data, &foo)
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	//fmt.Printf("%#v\n", foo)


	wand := imagick.NewMagickWand()
	pixel := imagick.NewPixelWand()
	pixel.SetColor("transparent")
	wand.NewImage(2000, 2000, pixel)
	draw := imagick.NewDrawingWand()
	err = draw.SetVectorGraphics(string(data))
	//draw.Annotation(100, 100, "foobar")
	//
	if err != nil {
		fmt.Println(err)
	}
	err = wand.DrawImage(draw)
	wand.Clear()

	err = wand.ReadImage("wandkalender_a3-hoch_month.svg")
	if err != nil {
		fmt.Println(err)
	}
	err = wand.WriteImage("foo.jpg")
	err = wand.WriteImage("foo.png")
	if err != nil {
		fmt.Println(err)
	}
}
