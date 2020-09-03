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
		CalendarType{},
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
