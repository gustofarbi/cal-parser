package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
)

func main() {
	str := "foobar and herer"
	x, y := 500, 500
	font := "AmaticSC-Regular.ttf"
	canvas := gg.NewContext(1000, 1000)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()
	err := canvas.LoadFontFace(font, 60)
	if err != nil {
		fmt.Println("font not found: " + font)
	}
	w, h := canvas.MeasureString(str)
	var r float64
	if w > h {
		r = w
	} else {
		r = h
	}
	text := gg.NewContext(int(r*2), int(r*2))
	err = text.LoadFontFace(font, 60)
	if err != nil {
		fmt.Println(err)
	}
	text.SetRGB(0, 0, 0)
	text.DrawString(str, r, r)
	text.RotateAbout(gg.Radians(-90), r, r)
	text.DrawString(str, r, r)
	text.RotateAbout(gg.Radians(180), r, r)
	text.DrawString(str, r, r)
	text.RotateAbout(gg.Radians(90), r, r)
	text.DrawString(str, r, r)
	text.SavePNG("text.png")
	canvas.DrawImageAnchored(text.Image(), x, y, 0.5, 0.5)
	canvas.SavePNG("foo.png")
}
