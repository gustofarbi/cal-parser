package svg

import (
	"regexp"
	"strconv"
)

type LineSkipDay struct {
	Attribute
}

func (l LineSkipDay) Apply(text *CalendarText) {}

func (l LineSkipDay) Matches(subject string) bool {
	reg := regexp.MustCompile("^n(\\d{1,2})")
	return reg.MatchString(subject)
}

func (l LineSkipDay) New(subject string) Annotation {
	reg := regexp.MustCompile("^n(\\d{1,2})")
	val, _ := strconv.Atoi(reg.FindStringSubmatch(subject)[1])
	return LineSkipDay{Attribute{val}}
}

func (l LineSkipDay) Id() string {
	return "lineSkipDay"
}

func (l LineSkipDay) Priority() int {
	return 0
}

