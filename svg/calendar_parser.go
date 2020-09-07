package svg

import (
	"regexp"
	"strings"
)

func (c Calendar) Parse(svg Svg, svgRaw string, scalingRatio float64) {
	c.svgContent = svgRaw
	go c.StartReceiver()
	context := NewContext(c.Receiver)
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
}

func parseGroup(g Group, ctx Context) {
	ctx = ctx.Merge(g)

	if ctx.RenderPrevNext() {
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
	calendarText := CalendarText{
		Position:   text.Position,
		Content:    text.Content,
		FontSize:   text.FontSize,
		FontFamily: text.FontFamily,
		FontColor:  text.Fill,
	}

	text.Tranform.Apply(calendarText)
	ctx.ApplyEarly(calendarText)
	calendarText.Annotations = ctx.Annotations

	ctx.Receiver <- calendarText
}

func (c Calendar) RemoveTexts() {
	reg := regexp.MustCompile("<text.*?/text>")
	reg.ReplaceAllString(c.svgContent, "")
}
