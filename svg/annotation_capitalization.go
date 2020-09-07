package svg

import (
	"strings"
)

type Capitalization struct {
	Attribute
}

func (c Capitalization) Apply(text *CalendarText) {
	switch c.Attr().(string) {
	case "upper":
		text.Content = strings.ToUpper(text.Content)
		break
	case "lower":
		text.Content = strings.ToLower(text.Content)
		break
	}
}

func (c Capitalization) Matches(subject string) bool {
	return subject == "upper" || subject == "lower"
}

func (c Capitalization) New(subject string) Annotation {
	return Capitalization{Attribute{subject}}
}

func (c Capitalization) Id() string {
	return "capitalization"
}

func (c Capitalization) Priority() int {
	return 99
}
