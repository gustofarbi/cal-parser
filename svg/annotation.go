package svg

type Annotation interface {
	Apply(text *CalendarText)
	Matches(subject string) bool
	New(subject string) Annotation
	Id() string
	Priority() int
	Attr() interface{}
}

var registered []Annotation

func init() {
	registered = []Annotation{
		Alignment{},
		CalendarType{},
		CalendarWeek{},
		Capitalization{},
		FormatMonthNumber{},
		FormatMonthText{},
		FormatWeekdayHeader{},
		FormatWeekdayPosition{},
		FormatYear{},
		HasRefinement{},
		Language{},
		LineSkipDay{},
		LineWeekdayElement{},
		LineWeekendElement{},
		Month{},
		RefinementType{},
		RenderMonthOnly{},
		Scaling{},
		DayAnotherMonth{},
		DayCurrentMonth{},
		RenderPrevNextMonth{},
		SkipWeek{},
		WeekdayHeader{},
		WeekdayPosition{},
	}
}

func Parse(object HasAnnotations) []Annotation {
	parsed := make([]Annotation, 0)
	for _, subject := range object.All() {
		for _, annotation := range registered {
			if annotation.Matches(subject) {
				parsed = append(parsed, annotation.New(subject))
			}
		}
	}

	return parsed
}
