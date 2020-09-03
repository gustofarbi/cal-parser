package svg

const (
	table = "table"
	line  = "line"
)

type CalendarType struct {
	Attr     string
}

func (c CalendarType) Apply(text CalendarText) {
	text.CalendarType = c.Attr
}

func (c CalendarType) Matches(subject string) bool {
	return subject == table || subject == line
}

func (c CalendarType) New(subject string) Annotation {
	return CalendarType{Attr: subject}
}

func (c CalendarType) Id() string {
	return "calendarType"
}

func (c CalendarType) Priority() int {
	return -1
}
