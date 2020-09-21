package svg

import "github.com/fogleman/gg"

type Alignment struct {
	Attribute
}

func (a Alignment) Apply(text *CalendarText) {
	ctx := gg.NewContext(0, 0)
	face, _ := GetFont(text.FontFamily, text.FontSize)
	ctx.SetFontFace(*face)
	w, _ := ctx.MeasureString(text.Content)
	switch a.Attr().(string) {
	case "c":
		text.Position.X += (text.Position.Width - w) / 2
	case "r":
		text.Position.X += text.Position.Width - w
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
