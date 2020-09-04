package svg

import (
	"regexp"
)

type FormatMonthNumber struct {
	format string
}

func (f FormatMonthNumber) Apply(text CalendarText) {
	if !text.IsMonth {
		return
	}

	// todo
	switch f.format {
	case "02":
	case "2":
	}
}

func (f FormatMonthNumber) Matches(subject string) bool {
	reg := regexp.MustCompile("^mn([02]{1,2})")
	return reg.MatchString(subject)
}

func (f FormatMonthNumber) New(subject string) Annotation {
	reg := regexp.MustCompile("^mn([02]{1,2})")
	val := reg.FindString(subject)
	return FormatMonthNumber{val}
}

func (f FormatMonthNumber) Id() string {
	return "formatMonthNumber"
}

func (f FormatMonthNumber) Priority() int {
	return 20
}

