package svg

import (
	"regexp"
	"strconv"
)

type FormatMonthNumber struct {
	Attribute
}

func (f FormatMonthNumber) Apply(text *CalendarText) {
	if !text.IsMonth {
		return
	}

	text.Content = strconv.Itoa(text.CurrentMonth)
	if f.Attr().(string) == "02" && len(text.Content) == 1 {
		text.Content = "0" + text.Content
	}
}

func (f FormatMonthNumber) Matches(subject string) bool {
	reg := regexp.MustCompile("^mn([02]{1,2})")
	return reg.MatchString(subject)
}

func (f FormatMonthNumber) New(subject string) Annotation {
	reg := regexp.MustCompile("^mn([02]{1,2})")
	val := reg.FindStringSubmatch(subject)
	return FormatMonthNumber{Attribute{val[1]}}
}

func (f FormatMonthNumber) Id() string {
	return "formatMonthNumber"
}

func (f FormatMonthNumber) Priority() int {
	return 20
}
