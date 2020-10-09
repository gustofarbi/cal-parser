package main

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
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
		DrawString(text, str, r, r+150, -0.1, &fontFace)
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

func DrawString(c *gg.Context, s string, x, y, letterSpacing float64, ff *font.Face) {
	formerX := x
	h := float64((*ff).Metrics().Height >> 6)
	for i, line := range strings.Split(s, "\n") {
		var prevC = rune(-1)
		y += float64(i) * 1.37 * h
		x = formerX
		for _, r := range line {
			var a fixed.Int26_6
			if prevC >= 0 {
				a += (*ff).Kern(prevC, r)
			}
			tmp, ok := (*ff).GlyphAdvance(r)
			if !ok {
				continue
			}
			c.DrawString(string(r), x, y)
			a += tmp
			af := float64(a >> 6)
			af += af * letterSpacing
			x += af
			prevC = r
		}
	}
}
