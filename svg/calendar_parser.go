package svg

import (
	"regexp"
	"strings"
)

func (c Calendar) Parse(svg Svg, scalingRatio float64) {
	go c.StartReceiver()
	context := Context{} // todo
	context.Add([]Annotation{
		Language{Attribute{"de"}},
		Alignment{Attribute{"r"}},
		Scaling{Attribute{scalingRatio}},
		FormatWeekdayHeader{Attribute{"dt2"}},
		FormatWeekdayPosition{Attribute{"nn2"}},
	})

	for _, g := range svg.Gs {
		if strings.Contains(g.Id, "calendar") {
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
		RenderMonthOnly{},
		SkipWeek{},
		LineSkipDay{},
		LineWeekdayElement{},
		LineWeekendElement{}, // todo: weekday position is missing
	}, g.Raw)

	for _, text := range g.Texts {
		go parseText(text, ctx)
	}
	for _, group := range g.Gs {
		go parseGroup(group, ctx) // todo maybe not necessary
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
