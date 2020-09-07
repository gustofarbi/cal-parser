package svg

type LineWeekdayElement struct {
	Attribute
}

func (l LineWeekdayElement) Apply(text CalendarText) {
	text.IsWeekday = true
}

func (l LineWeekdayElement) Matches(subject string) bool {
	return subject == "wd"
}

func (l LineWeekdayElement) New(subject string) Annotation {
	return LineWeekdayElement{}
}

func (l LineWeekdayElement) Id() string {
	return "lineWeekdayElement"
}

func (l LineWeekdayElement) Priority() int {
	return -1
}
