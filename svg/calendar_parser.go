package svg

import (
	"fmt"
	"github.com/fogleman/gg"
	"regexp"
	"strings"
)

func (c *Calendar) Parse(svg Svg, svgRaw string, scalingRatio float64) {
	c.svgContent = svgRaw
	go c.StartReceiver()
	context := NewContext(c.Receiver, c.ReceiverWg)
	context.Add([]Annotation{
		Language{Attribute{"de"}},
		Alignment{Attribute{"r"}},
		Scaling{Attribute{scalingRatio}},
		FormatWeekdayHeader{Attribute{"2"}},
		FormatWeekdayPosition{Attribute{"2"}},
	})

	for _, g := range svg.Gs {
		if strings.Contains(g.Id, "month") || strings.Contains(g.Id, "calendar") {
			parseGroup(g, context)
		}
	}
	c.RemoveTexts()
	//c.ReceiverWg.Wait()
}

func parseGroup(g Group, formerCtx Context) {
	ctx := formerCtx.Merge(g)

	if ctx.RenderPrevNext() {
		ctx.Receiver <- AnnotationObject{RenderPrevNextMonth{Attribute{true}}, ""}
	}

	ctx.HandleSpecialAnnotation([]Annotation{
		SkipWeek{},
		LineSkipDay{},
		LineWeekdayElement{},
		LineWeekendElement{},
	}, g.Identifier())

	for _, text := range g.Texts {
		ctx.ReceiverWg.Add(1)
		parseText(text, ctx)
	}
	for _, rect := range g.Rects {
		parseDataName(rect, ctx)
	}
	for _, circle := range g.Circles {
		parseDataName(circle, ctx)
	}
	for _, group := range g.Gs {
		parseGroup(group, ctx)
	}
}

func parseDataName(object HasAnnotations, formerCtx Context) {
	ctx := formerCtx.Merge(object)

	if ctx.RenderPrevNext() {
		ctx.Receiver <- AnnotationObject{RenderPrevNextMonth{Attribute{true}}, ""}
	}

	ctx.HandleSpecialAnnotation([]Annotation{
		SkipWeek{},
		LineSkipDay{},
		LineWeekdayElement{},
		LineWeekendElement{},
	}, object.Identifier())
}

func parseText(text Text, ctx Context) {
	ctx = ctx.Merge(text)
	calendarText := CalendarText{
		Position:   text.Position,
		Content:    text.Content,
		FontSize:   text.FontSize,
		FontFamily: text.FontFamily,
		FontColor:  text.Fill,
	}

	calendarText.Position.Width, calendarText.Position.Height = calculateDimensions(
		calendarText.Content, calendarText.FontFamily, calendarText.FontSize)
	text.Tranform.Apply(&calendarText)
	ctx.Annotations.ApplyEarly(&calendarText)
	calendarText.Annotations = ctx.Annotations

	ctx.Receiver <- calendarText
}

func calculateDimensions(text, fontString string, fontSize float64) (float64, float64) {
	ctx := gg.NewContext(0, 0)
	face, err := GetFont(fontString, fontSize)
	if err != nil {
		panic(fmt.Sprintf("font-string could not be parsed: %s", fontString))
	}
	ctx.SetFontFace(*face)
	return ctx.MeasureString(text)
}

func (c *Calendar) RemoveTexts() {
	reg := regexp.MustCompile("<text.*?/text>")
	c.svgContent = reg.ReplaceAllString(c.svgContent, "")
}
