package svg

type DayAnotherMonth struct {
	Attribute
}

func (d DayAnotherMonth) Apply(text *CalendarText) {
	text.IsAnotherMonth = true
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
