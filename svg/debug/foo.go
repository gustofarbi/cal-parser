package main

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image/color"
	"io/ioutil"
	"strings"
)

func main() {
	str := "foobar and here\nplus this here"
	x, y := 500, 500
	fontfile := "resources/fonts/ArchitectsDaughter-Regular.ttf"
	canvas := gg.NewContext(1000, 1000)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()
	fontData, err := ioutil.ReadFile(fontfile)
	if err != nil {
		panic("fontfile not found: " + fontfile + " " + err.Error())
	}
	fontRes, err := truetype.Parse(fontData)
	//c := freetype.NewContext()
	//c.SetFont(fontRes)
	//c.
	if err != nil {
		panic("fontRes cannot be parsed: " + err.Error())
	}
	fontSize := 60.0
	fontFace := truetype.NewFace(fontRes, &truetype.Options{
		Size:       fontSize,
		SubPixelsX: 10,
		Hinting:    font.HintingFull,
	})
	canvas.SetFontFace(fontFace)
	w, h := canvas.MeasureString(str)
	var r float64
	if w > h {
		r = w
	} else {
		r = h
	}
	text := gg.NewContext(int(r*2), int(r*2))
	text.SetFontFace(fontFace)
	text.SetRGB(0, 0, 0)
	for i := 0; i < 4; i++ {
		for j, line := range strings.Split(str, "\n") {
			text.DrawString(line, r, r+float64(j)*1.37*h)
		}
		DrawString(text, str, r, r, 1.0, &fontFace)
		text.RotateAbout(gg.Radians(90), r, r)
	}
	err = text.SavePNG("text.png")
	if err != nil {
		panic("cannot save text png")
	}
	canvas.DrawImageAnchored(text.Image(), x, y, 0.5, 0.5)
	err = canvas.SavePNG("foo.png")
	if err != nil {
		panic("cannot save canvas png")
	}
}

type ContextWithLetterSpacing struct {
	gg.Context
	letterSpacing float64
}

func NewContext(width, height int, letterSpacing float64) *ContextWithLetterSpacing {
	ctx := gg.NewContext(width, height)
	return &ContextWithLetterSpacing{
		*ctx,
		letterSpacing,
	}
}

func DrawString(c *gg.Context, s string, x, y, letterSpacing float64, ff *font.Face) {
	formerX := x
	_, h := c.MeasureString(s)
	for i, line := range strings.Split(s, "\n") {
		y += float64(i) * 1.37 * h
		x = formerX
		for j, r := range line {
			if j != 0 {
				a, ok := (*ff).GlyphAdvance(r)
				if !ok {
					panic("not OK")
				}
				af := 0.0
				af += float64(a.Ceil()) + float64(a.Floor()/1000000)
				x += af
			}
			c.DrawString(string(r), x, y)
		}
	}
}
