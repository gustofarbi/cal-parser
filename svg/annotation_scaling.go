package svg

type Scaling struct {
	Attribute
}

func (s Scaling) Apply(text *CalendarText) {
	factor := s.Attr().(float64)
	x := text.Position.X
	y := text.Position.Y
	width := text.Position.Width
	height := text.Position.Height
	fz := text.FontSize
	text.Position.X = x * factor
	text.Position.Y = y * factor
	text.Position.Width = width * factor
	text.Position.Height = height * factor
	text.FontSize = fz * factor
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
