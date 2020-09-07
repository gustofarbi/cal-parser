package svg

import (
	"regexp"
	"strconv"
)

type WeekdayHeader struct {
	Attribute
}

func (w WeekdayHeader) Apply(text *CalendarText) {
	text.WeekdayHeader = w.Attr().(int)
}

func (w WeekdayHeader) Matches(subject string) bool {
	reg := regexp.MustCompile("^d(\\d{1,2})")
	return reg.MatchString(subject)
}

func (w WeekdayHeader) New(subject string) Annotation {
	reg := regexp.MustCompile("^d(\\d{1,2})")
	raw := reg.FindStringSubmatch(subject)
	val, _ := strconv.Atoi(raw[1])
	return WeekdayHeader{Attribute{val}}
}

func (w WeekdayHeader) Id() string {
	return "weekdayHeader"
}

func (w WeekdayHeader) Priority() int {
	return -1
}

