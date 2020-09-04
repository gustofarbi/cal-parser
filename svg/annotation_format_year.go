package svg

import "regexp"

type FormatYear struct {
	format string
}

func (f FormatYear) Apply(text CalendarText) {
	if !text.IsYear {
		return
	}

	if f.format == "2" {
		text.Content = text.Content[len(text.Content)-2:]
	}
}

func (f FormatYear) Matches(subject string) bool {
	reg := regexp.MustCompile("^yn([24])")
	return reg.MatchString(subject)
}

func (f FormatYear) New(subject string) Annotation {
	reg := regexp.MustCompile("^yn([24])")
	val := reg.FindString(subject)
	return FormatYear{val}
}

func (f FormatYear) Id() string {
	return "formatYear"
}

func (f FormatYear) Priority() int {
	return 20
}
