package svg

import (
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var weekdays = map[int]string{
	1: "Montag",
	2: "Dienstag",
	3: "Mittwoch",
	4: "Donnerstag",
	5: "Freitag",
	6: "Samstag",
	7: "Sonntag",
}

var (
	ch = make(chan ImageObject)
	wg = &sync.WaitGroup{}

)

func (c *Calendar) Render(year, month int, width float64) {
	ctx := renderSvg(c.svgContent, width)
	go startReceiver(ctx)
	c.drawTexts(year, month)
	ctx.SavePNG("saved.png")
}

type ImageObject struct {
	im *gg.Context
	p  image.Point
}

func startReceiver(canvas *gg.Context) {
	for o := range ch {
		fmt.Println("drawing text")
		canvas.DrawImageAnchored(o.im.Image(), o.p.X, o.p.X, 0.5, 0.5)
		wg.Done()
	}
}

func (c *Calendar) drawTexts(year, month int) {
	if len(c.positionTableCurrentMonth) == 42 {
		c.fillTable(year, month)
	}
	wg.Wait()
	close(ch)
}

func drawSingleText(c *CalendarText) {
	c.Annotations.ApplyLate(c)
	w, h := c.Position.Width, c.Position.Height
	var r float64
	if w > h {
		r = w
	} else {
		r = h
	}
	ctx := gg.NewContext(int(r*2), int(r*2))
	// todo fonts
	err := ctx.LoadFontFace("AmaticSC-Regular.ttf", c.FontSize)
	if err != nil {
		panic(err)
	}
	ctx.RotateAbout(gg.Radians(c.Position.Rotation), r, r)
	var fontColor color.Color
	if c.FontColor != "" {
		fontColor, err = ParseHexColor(c.FontColor)
		if err != nil {
			panic(err)
		}
	} else {
		fontColor = color.Black
	}
	ctx.SetColor(fontColor)
	ctx.DrawString(c.Content, r, r)

	ch <- ImageObject{
		ctx, image.Point{
			X: int(c.Position.X),
			Y: int(c.Position.Y),
		},
	}
}

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}

func (c *Calendar) fillTable(year, month int) {
	for _, w := range c.weekdayHeadingsTable {
		w.Content = weekdays[w.WeekdayHeader]
		wg.Add(1)
		go drawSingleText(&w)
	}

	//for pos := range c.positionTableCurrentMonth {
	//
	//}
}

func renderSvg(svg string, width float64) *gg.Context {
	svgFile, _ := ioutil.TempFile("", "*.svg")
	defer os.Remove(svgFile.Name())
	pngFile, _ := ioutil.TempFile("", "*.png")
	defer os.Remove(pngFile.Name())

	_, err := svgFile.WriteString(svg)
	if err != nil {
		fmt.Println(err)
	}


	cmd := exec.Command("/usr/local/bin/rsvg-convert", "-w", strconv.Itoa(int(width)), svgFile.Name(), "-o", pngFile.Name())
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	im, _, err := image.Decode(pngFile)
	if err != nil {
		fmt.Println(err)
	}

	return gg.NewContextForImage(im)
}
