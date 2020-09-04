package svg

import "regexp"

type LineSkipDay struct {}

func (l LineSkipDay) Apply(text CalendarText) {}

func (l LineSkipDay) Matches(subject string) bool {
	reg := regexp.MustCompile("^n(\\d{1,2})")
	return reg.MatchString(subject)
}

func (l LineSkipDay) New(subject string) Annotation {
	return LineSkipDay{}
}

func (l LineSkipDay) Id() string {
	return "lineSkipDay"
}

func (l LineSkipDay) Priority() int {
	return 0
}

