package svg

import (
	"regexp"
	"strconv"
)

type WeekdayHeader struct {
	pos int
}

func (w WeekdayHeader) Apply(text CalendarText) {
	text.WeekdayHeader = w.pos
}

func (w WeekdayHeader) Matches(subject string) bool {
	reg := regexp.MustCompile("^d(\\d{1,2})")
	return reg.MatchString(subject)
}

func (w WeekdayHeader) New(subject string) Annotation {
	val, _ := strconv.Atoi(subject)
	return WeekdayHeader{val}
}

func (w WeekdayHeader) Id() string {
	return "weekdayHeader"
}

func (w WeekdayHeader) Priority() int {
	return -1
}

