package svg

type LineWeekendElement struct {
	Attribute
}

func (l LineWeekendElement) Apply(text *CalendarText) {
	text.IsWeekend = true
}

func (l LineWeekendElement) Matches(subject string) bool {
	return subject == "we"
}

func (l LineWeekendElement) New(subject string) Annotation {
	return LineWeekendElement{}
}

func (l LineWeekendElement) Id() string {
	return "lineWeekendElement"
}

func (l LineWeekendElement) Priority() int {
	return -1
}

