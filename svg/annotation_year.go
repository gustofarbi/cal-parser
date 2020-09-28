package svg

import "strings"

type Year struct {
	Attribute
}

func (y Year) Apply(text *CalendarText) {
	text.IsYear = true
}

func (y Year) Matches(subject string) bool {
	return strings.HasPrefix(subject, "yn")
}

func (y Year) New(subject string) Annotation {
	return Year{} // todo
}

func (y Year) Id() string {
	return "year"
}

func (y Year) Priority() int {
	return -1
}
