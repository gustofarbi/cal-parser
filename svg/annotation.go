package svg

type Annotation interface {
	Apply(text CalendarText)
	Matches(subject string) bool
	New(subject string) Annotation
	Id() string
	Priority() int
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

func Parse(object HasAnnotations) map[string]Annotation {
	parsed := make(map[string]Annotation)
	for _, subject := range object.All() {
		for _, annotation := range registered {
			if annotation.Matches(subject) {
				parsed[annotation.Id()] = annotation.New(subject)
			}
		}
	}

	return parsed
}
