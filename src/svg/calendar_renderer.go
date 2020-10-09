package svg

import (
	"fmt"
	"github.com/fogleman/gg"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var weekdays = map[int]string{
	1: "Montag",
	2: "Dienstag",
	3: "Mittwoch",
	4: "Donnerstag",
	5: "Freitag",
	6: "Samstag",
	0: "Sonntag",
}

var (
	ch = make(chan ImageObject)
	wg = &sync.WaitGroup{}
)

func (c *Calendar) Render(year, month int, width, height float64) {
	canvas := gg.NewContext(int(width), int(height))
	go startReceiver(canvas)
	c.drawTexts(year, month)
	bg := renderSvg(c.svgContent, width)
	bg.DrawImage(canvas.Image(), 0, 0)
	bg.SavePNG("saved.png")
}

type ImageObject struct {
	im *gg.Context
	p  image.Point
}

func startReceiver(canvas *gg.Context) {
	for o := range ch {
		canvas.DrawImageAnchored(o.im.Image(), o.p.X, o.p.Y, 0.5, 0.5)
		wg.Done()
	}
}

func FirstOfTheMonth(year, month int) time.Time {
	return time.Date(year, time.Month(month), 1, 12, 0, 0, 0, time.Local)
}

func (c *Calendar) drawTexts(year, month int) {
	if len(c.positionTableCurrentMonth) == 42 {
		c.fillTable(year, month)
	}

	if len(c.positionsLineDefault) == 31 || len(c.positionsLineWeekday) == 31 {
		c.fillLine(year, month)
	}

	for _, text := range c.months {
		wg.Add(1)
		drawSingleText(&text, year, month)
	}

	for _, text := range c.years {
		text.Content = strconv.Itoa(FirstOfTheMonth(year, month).Year())
		wg.Add(1)
		drawSingleText(&text, year, month)
	}

	for _, text := range c.texts {
		wg.Add(1)
		drawSingleText(&text, year, month)
	}
	wg.Wait()
	close(ch)
}

func drawSingleText(c *CalendarText, year, month int) {
	c.CurrentMonth = month
	c.CurrentYear = year
	c.Annotations.ApplyLate(c)
	ctx := gg.NewContext(0, 0)
	face, err := GetFont(c.FontFamily, c.FontSize)
	if err != nil {
		panic(err)
	}
	ctx.SetFontFace(*face)
	w, h := ctx.MeasureString(c.Content)
	var r float64
	if w > h {
		r = w
	} else {
		r = h
	}
	r = math.Sqrt(w*w + h*h)
	ctx = gg.NewContext(int(r*2), int(r*2))
	ctx.SetFontFace(*face)
	ctx.RotateAbout(gg.Radians(c.Position.Rotation), r, r)
	var fontColor color.Color
	var ok bool
	if c.FontColor != "" {
		if strings.HasPrefix(c.FontColor, "#") {
			fontColor, err = ParseHexColor(c.FontColor)
			if err != nil {
				panic(err)
			}
		} else {
			fontColor, ok = colornames.Map[c.FontColor]
			if !ok {
				panic("color not found: " + c.FontColor)
			}
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
	// render headers first
	for _, w := range c.weekdayHeadingsTable {
		w.Content = weekdays[w.WeekdayHeader%7]
		wg.Add(1)
		drawSingleText(&w, year, month)
	}

	// render days in previous month if necessary
	counter := 0
	currentDate := FirstOfTheMonth(year, month)
	startWeekday := currentDate.Weekday()
	if startWeekday != time.Monday && c.RenderPrevNext {
		for i := int(startWeekday); i > 0; i-- {
			text, ok := c.positionTableAnotherMonth[i%7]
			if !ok {
				continue
			}
			text.Content = strconv.Itoa(currentDate.Day())
			wg.Add(1)
			drawSingleText(&text, year, month)
			counter++
			currentDate = currentDate.AddDate(0, 0, -1)
		}
	}

	// render current month
	currentDate = FirstOfTheMonth(year, month)
loop:
	for i := 0; i < 7; i++ {
		for d := int(startWeekday); d < 7; d++ {
			if int(currentDate.Month()) == month {
				text := c.positionTableCurrentMonth[i*7+d]
				text.Content = strconv.Itoa(currentDate.Day())
				wg.Add(1)
				drawSingleText(&text, year, month)
				counter++
				currentDate = currentDate.AddDate(0, 0, 1)
				startWeekday = 0
			} else {
				break loop
			}
		}
	}

	// render days in the month after if necessary
	if c.RenderPrevNext && currentDate.Weekday() != time.Monday {
		var text CalendarText
		var ok bool
		for currentDate.Weekday() != time.Monday {
			text, ok = c.positionTableAnotherMonth[counter]
			if !ok {
				text, ok = c.positionTableCurrentMonth[counter]
			}
			text.Content = strconv.Itoa(currentDate.Day())
			wg.Add(1)
			drawSingleText(&text, year, month)
			counter++
			currentDate = currentDate.AddDate(0, 0, 1)
		}
	}

	currentDate = FirstOfTheMonth(year, month)
	for _, calendarWeek := range c.calendarWeeksTable {
		_, w := currentDate.ISOWeek()
		calendarWeek.Content = strconv.Itoa(w)
		wg.Add(1)
		drawSingleText(&calendarWeek, year, month)
		currentDate = currentDate.Add(7 * 24 * time.Hour)
		if int(currentDate.Month()) != month {
			break
		}
	}

	c.removeSkipDays(counter, true)
	c.removeSkipMonths(counter)
}

func (c *Calendar) removeSkipMonths(start int) {
	monthsToRemove := start % 7
	for i := 0; i < monthsToRemove; i++ {

	}
}

func (c *Calendar) removeSkipDays(start int, greaterThan bool) {
	for _, node := range c.nodeMapping.skipDays {
		for _, skipDay := range node.daysToSkip {
			if greaterThan && skipDay >= start {
				c.removeNode(node)
			} else if !greaterThan && skipDay <= start {
				c.removeNode(node)
			}
		}
	}
}

func (c *Calendar) removeNode(n *Node) {
	for _, child := range n.Children {
		c.removeNode(child)
	}

	if n.end == "" {
		c.svgContent = strings.Replace(c.svgContent, n.start, "", -1)
	} else {
		c.svgContent = strings.Replace(c.svgContent, n.start+n.end, "", -1) // todo: check if not regex
	}
}

func (c *Calendar) fillLine(year, month int) {
	first := time.Date(year, time.Month(month), 1, 12, 0, 0, 0, time.Local)
	counter := 1
	for int(first.Month()) == month {
		header, ok := c.weekdayHeadingsLine[counter]
		if !ok {
			panic("header not set: " + strconv.Itoa(counter))
		}
		header.Content = weekdays[int(first.Weekday())%7]
		first = first.Add(24 * time.Hour)
		counter++
		wg.Add(1) // todo do this in the method itself
		drawSingleText(&header, year, month)
	}
	first = time.Date(year, time.Month(month), 1, 12, 0, 0, 0, time.Local)
	counter = 1
	for int(first.Month()) == month {
		isWeekend := first.Weekday() == time.Saturday || first.Weekday() == time.Sunday
		c.fillPositionLine(counter, year, month, isWeekend)
		counter++
		first = first.Add(24 * time.Hour)
	}
	c.removeSkipDays(counter, true)
}

func (c *Calendar) fillPositionLine(position, year, month int, isWeekend bool) {
	var pos CalendarText
	var ok bool
	if isWeekend {
		pos, ok = c.positionsLineWeekend[position]
	} else {
		pos, ok = c.positionsLineWeekday[position]
	}
	if !ok {
		pos, ok = c.positionsLineDefault[position]
		if !ok {
			panic("no position found: " + strconv.Itoa(position))
		}
	}
	wg.Add(1)
	drawSingleText(&pos, year, month)
}

func renderSvg(svg string, width float64) *gg.Context {
	svgFile, _ := ioutil.TempFile("", "*.src")
	defer os.Remove(svgFile.Name())
	pngFile, _ := ioutil.TempFile("", "*.png")
	defer os.Remove(pngFile.Name())

	_, err := svgFile.WriteString(svg)
	if err != nil {
		fmt.Println(err)
	}

	cmd := exec.Command("rsvg-convert", "-w", strconv.Itoa(int(width)), svgFile.Name(), "-o", pngFile.Name())
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
