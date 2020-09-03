package svg

type Scaling struct {
	scale float64
}

func (s Scaling) Apply(text CalendarText) {
	text.Position.X *= s.scale
	text.Position.Y *= s.scale
	text.FontSize *= s.scale
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
