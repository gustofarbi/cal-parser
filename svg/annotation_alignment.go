package svg


type Alignment struct {
	Attribute
}

func (a Alignment) Apply(text *CalendarText) {
	switch a.Attr().(string) {
	case "c":
		//text.Position.X += (text.Position.Width - dims.TextWidth) / 2
	case "r":
		//text.Position.X += text.Position.Width - dims.TextWidth
	}
}

func (a Alignment) Matches(subject string) bool {
	return subject == "r" || subject == "c" || subject == "l"
}

func (a Alignment) New(subject string) Annotation {
	return Alignment{Attribute{subject}}
}

func (a Alignment) Id() string {
	return "alignment"
}

func (a Alignment) Priority() int {
	return 100
}
