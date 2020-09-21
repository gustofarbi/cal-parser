package svg

type DayCurrentMonth struct {
	Attribute
}

func (d DayCurrentMonth) Apply(text *CalendarText) {
	text.IsAnotherMonth = false
}

func (d DayCurrentMonth) Matches(subject string) bool {
	return subject == "mc"
}

func (d DayCurrentMonth) New(subject string) Annotation {
	return DayCurrentMonth{}
}

func (d DayCurrentMonth) Id() string {
	return "dayCurrentMonth"
}

func (d DayCurrentMonth) Priority() int {
	return -1
}
