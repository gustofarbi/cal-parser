package svg

import "strings"

type Month struct {}

func (m Month) Apply(text CalendarText) {
	text.IsMonth = true
}

func (m Month) Matches(subject string) bool {
	return strings.HasPrefix(subject, "mn") || strings.HasPrefix(subject, "mt")
}

func (m Month) New(subject string) Annotation {
	return Month{}
}

func (m Month) Id() string {
	return "month"
}

func (m Month) Priority() int {
	return -1
}
