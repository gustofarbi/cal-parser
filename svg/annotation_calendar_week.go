package svg

import (
	"regexp"
	"strconv"
)

type CalendarWeek struct {
	Attribute
}

func (c CalendarWeek) Apply(text *CalendarText) {
	val, _ := strconv.Atoi(c.Attr().(string))
	text.CalendarWeek = val
}

func (c CalendarWeek) Matches(subject string) bool {
	reg := regexp.MustCompile("^kw(\\d{1,2})")
	return reg.MatchString(subject)
}

func (c CalendarWeek) New(subject string) Annotation {
	reg := regexp.MustCompile("^kw(\\d{1,2})")
	val := reg.FindStringSubmatch(subject)
	return CalendarWeek{Attribute{val[1]}}
}

func (c CalendarWeek) Id() string {
	return "calendarWeek"
}

func (c CalendarWeek) Priority() int {
	return -1
}
