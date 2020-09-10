package svg

type Scaling struct {
	Attribute
}

func (s Scaling) Apply(text *CalendarText) {
	text.Position.X = text.Position.X * s.Attr().(float64)
	text.Position.Y = text.Position.Y * s.Attr().(float64)
	text.FontSize = text.FontSize * s.Attr().(float64)
}

func (s Scaling) Matches(subject string) bool {
	return false
}

func (s Scaling) New(subject string) Annotation {
	return Scaling{}
}

func (s Scaling) Id() string {
	return "scaling"
}

func (s Scaling) Priority() int {
	return 101
}
