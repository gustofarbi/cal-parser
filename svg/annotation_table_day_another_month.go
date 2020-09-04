package svg

type DayAnotherMonth struct {}

func (d DayAnotherMonth) Apply(text CalendarText) {
	text.IsCurrentMonth = false
}

func (d DayAnotherMonth) Matches(subject string) bool {
	return subject == "mpn"
}

func (d DayAnotherMonth) New(subject string) Annotation {
	return DayAnotherMonth{}
}

func (d DayAnotherMonth) Id() string {
	return "dayAnotherMonth"
}

func (d DayAnotherMonth) Priority() int {
	return -1
}

