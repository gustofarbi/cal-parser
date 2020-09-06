package svg

import (
	"regexp"
	"strconv"
)

type RenderMonthOnly struct {
	Attribute
}

func (r RenderMonthOnly) Apply(text CalendarText) {}

func (r RenderMonthOnly) Matches(subject string) bool {
	reg := regexp.MustCompile("^m(\\d{1,2})")
	return reg.MatchString(subject)
}

func (r RenderMonthOnly) New(subject string) Annotation {
	reg := regexp.MustCompile("^m(\\d{1,2})")
	val, _ := strconv.ParseBool(reg.FindString(subject))
	return RenderMonthOnly{Attribute{val}}
}

func (r RenderMonthOnly) Id() string {
	return "renderMonthOnly"
}

func (r RenderMonthOnly) Priority() int {
	return 0
}
