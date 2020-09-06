package svg

import (
	"strings"
)

func (c Calendar) Parse(svg Svg, scalingRatio float64) {
	go c.StartReceiver()
	c.Context = Context{} // todo
	c.Context.Add([]Annotation{
		Language{Attribute{"de"}},
		Alignment{Attribute{"r"}},
		Scaling{Attribute{scalingRatio}},
		FormatWeekdayHeader{Attribute{"dt2"}},
		FormatWeekdayPosition{Attribute{"nn2"}},
	})
	for _, g := range svg.Gs {
		if strings.Contains(g.Id, "calendar") {
			c.parseGroup(g)
		}
	}
}

func (c Calendar) parseGroup(g Group) {
	c.Context.Update(g)

	if !c.RenderPrevNext && c.Context.RenderPrevNext() {
		c.RenderPrevNext = true
	}

	c.Context.HandleSpecialAnnotation(RenderMonthOnly{})
	c.Context.HandleSpecialAnnotation(SkipWeek{})
	c.Context.HandleSpecialAnnotation(LineSkipDay{})
	c.Context.HandleSpecialAnnotation(WeekdayPosition{}) // todo: these need different kind of handling
	c.Context.HandleSpecialAnnotation(LineWeekendElement{})

	for _, text := range g.Texts {
		go c.parseText(text)
	}
	for _, group := range g.Gs {
		c.parseGroup(group)
	}
}

func (c Calendar) parseText(text Text) {
	calendarText := CalendarText{
		Position:   text.Position,
		Content:    text.Content,
		FontSize:   text.FontSize,
		FontFamily: text.FontFamily,
		FontColor:  text.Fill,
	}

	text.Tranform.Apply(calendarText)
	c.Context.ApplyEarly(calendarText)
	calendarText.Annotations = c.Context.Annotations

	c.Receiver <- calendarText
}
