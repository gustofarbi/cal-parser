package svg

import "regexp"

type FormatMonthText struct {
	Attribute
}

func (f FormatMonthText) Apply(text *CalendarText) {
	if !text.IsMonth {
		return
	}
	switch f.Attr().(string) {
	case "1":
		text.Content = text.Content[:1]
		break
	case "3":
		// todo
		text.Content = text.Content[:3]
		break
	}
}

func (f FormatMonthText) Matches(subject string) bool {
	reg := regexp.MustCompile("^mt([139])")
	return reg.MatchString(subject)
}

func (f FormatMonthText) New(subject string) Annotation {
	reg := regexp.MustCompile("^mt([139])")
	val := reg.FindStringSubmatch(subject)
	return FormatMonthText{Attribute{val[1]}}
}

func (f FormatMonthText) Id() string {
	return "formatMonthText"
}

func (f FormatMonthText) Priority() int {
	return 20
}
