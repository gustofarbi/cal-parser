package svg

import (
	"regexp"
	"strconv"
)

type WeekdayPosition struct {
	Attribute
}

func (w WeekdayPosition) Apply(text CalendarText) {
	text.WeekdayPosition = w.Attr().(int)
}

func (w WeekdayPosition) Matches(subject string) bool {
	reg := regexp.MustCompile("^n(\\d{1,2})")
	return reg.MatchString(subject)
}

func (w WeekdayPosition) New(subject string) Annotation {
	val, _ := strconv.Atoi(subject)
	return WeekdayPosition{Attribute{val}}
}

func (w WeekdayPosition) Id() string {
	return "weekdayPosition"
}

func (w WeekdayPosition) Priority() int {
	return -1
}

