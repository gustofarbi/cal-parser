package svg

import (
	"strings"
)

func (c Calendar) Parse(svg Svg) {
	c.Context = Context{make(map[string]Annotation)}
	for _, g := range svg.Gs {
		if strings.Contains(g.Id, "calendar") {
			c.parseCalendar(g)
		}
	}
}

func (c Calendar) parseGroup(g Group) {
	c.Context.Update(g)
	if strings.Contains(g.Id, "calendar") {
		c.parseCalendar(g)
	}
	for _, text := range g.Texts {
		c.parseText(text)
	}
	for _, group := range g.Gs {
		c.parseGroup(group)
	}
}

func (c Calendar) parseText(text Text) {
	//calendarText := CalendarText(text)
	calendarText := CalendarText{
		Position:   text.Position,
		Content:    text.Content,
		FontSize:   text.FontSize,
		FontFamily: text.FontFamily,
		FontColor:  text.Fill,
	}
	text.Tranform.Apply(calendarText)
	c.Context.ApplyEarly(calendarText)
	c.texts = append(c.texts, calendarText)
}

func (c Calendar) parseCalendar(g Group) {

}
