package svg

import (
	"fmt"
	"regexp"
	"strconv"
)

type FormatWeekdayPosition struct {
	Attribute
}

func (f FormatWeekdayPosition) Apply(text CalendarText) {
	if text.WeekdayPosition == 0 {
		return
	}
	if f.Attr().(string) == "02" {
		number, _ :=  strconv.Atoi(text.Content)
		text.Content = fmt.Sprintf("%02d", number)
	}
}

func (f FormatWeekdayPosition) Matches(subject string) bool {
	reg := regexp.MustCompile("^nn([02]{1,2})")
	return reg.MatchString(subject)
}

func (f FormatWeekdayPosition) New(subject string) Annotation {
	reg := regexp.MustCompile("^nn([02]{1,2})")
	val := reg.FindString(subject)
	return FormatWeekdayPosition{Attribute{val}}
}

func (f FormatWeekdayPosition) Id() string {
	return "formatWeekdayPosition"
}

func (f FormatWeekdayPosition) Priority() int {
	return 20
}

