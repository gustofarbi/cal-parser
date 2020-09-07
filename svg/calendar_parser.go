package svg

import (
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
		FormatWeekdayHeader{Attribute{"dt2"}},
		FormatWeekdayPosition{Attribute{"nn2"}},
	})

	for _, g := range svg.Gs {
		if strings.Contains(g.Id, "month") {
			parseGroup(g, context)
		}
	}
	c.RemoveTexts()
	c.ReceiverWg.Wait()
}

func parseGroup(g Group, formerCtx Context) {
	ctx := formerCtx.Merge(g.DataName)

	if ctx.RenderPrevNext() {
		ctx.ReceiverWg.Add(1)
		ctx.Receiver <- RenderPrevNextMonth{Attribute{true}}
	}

	ctx.HandleSpecialAnnotation([]Annotation{
		//RenderMonthOnly{}, // todo: other way
		SkipWeek{},
		LineSkipDay{},
		LineWeekdayElement{},
		LineWeekendElement{},
	}, g.Raw)

	for _, text := range g.Texts {
		parseText(text, ctx)
	}
	for _, group := range g.Gs {
		parseGroup(group, ctx) // todo maybe not necessary
	}
}

func parseText(text Text, ctx Context) {
	ctx = ctx.Merge(text.DataName)
	calendarText := CalendarText{
		Position:   text.Position,
		Content:    text.Content,
		FontSize:   text.FontSize,
		FontFamily: text.FontFamily,
		FontColor:  text.Fill,
	}

	text.Tranform.Apply(calendarText)
	ctx.ApplyEarly(&calendarText)
	calendarText.Annotations = ctx.Annotations

	ctx.ReceiverWg.Add(1)
	ctx.Receiver <- calendarText
}

func (c *Calendar) RemoveTexts() {
	reg := regexp.MustCompile("<text.*?/text>")
	reg.ReplaceAllString(c.svgContent, "")
}
