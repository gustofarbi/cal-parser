package svg

import "gopkg.in/gographics/imagick.v3/imagick"

type Alignment struct {
	Attribute
}

func (a Alignment) Apply(text *CalendarText) {
	wand := imagick.NewMagickWand()
	draw := imagick.NewDrawingWand()
	err := draw.SetFontFamily(text.FontFamily)
	draw.SetFontSize(text.FontSize)
	dims := wand.QueryFontMetrics(draw, text.Content)
	if err != nil {
		// todo
	}
	switch a.Attr().(string) {
	case "c":
		text.Position.X += (text.Position.Width - dims.TextWidth) / 2
		break
	case "r":
		text.Position.X += text.Position.Width - dims.TextWidth
		break
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
