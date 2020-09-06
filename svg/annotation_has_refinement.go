package svg

type HasRefinement struct {
	Attribute
}

func (h HasRefinement) Apply(text CalendarText) {
	text.HasRefinement = true
}

func (h HasRefinement) Matches(subject string) bool {
	return subject == "ref"
}

func (h HasRefinement) New(subject string) Annotation {
	return HasRefinement{}
}

func (h HasRefinement) Id() string {
	return "hasRefinement"
}

func (h HasRefinement) Priority() int {
	return 0
}
