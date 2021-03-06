package svg

import "regexp"

type FormatWeekdayHeader struct {
	Attribute
}

func (f FormatWeekdayHeader) Apply(text *CalendarText) {
	if text.WeekdayHeader == 0 {
		return
	}

	format := f.Attr().(string)
	switch format {
	case "1":
		text.Content = text.Content[:1]
	case "2":
		text.Content = text.Content[:2]
	case "2p":
		// todo: language
		text.Content = text.Content[:2]
		text.Content += "."
	}
}

func (f FormatWeekdayHeader) Matches(subject string) bool {
	reg := regexp.MustCompile("^dt([129p]{1,2})")
	return reg.MatchString(subject)
}

func (f FormatWeekdayHeader) New(subject string) Annotation {
	reg := regexp.MustCompile("^dt([129p]{1,2})")
	val := reg.FindStringSubmatch(subject)
	return FormatWeekdayHeader{Attribute{val[1]}}
}

func (f FormatWeekdayHeader) Id() string {
	return "formatWeekdayHeader"
}

func (f FormatWeekdayHeader) Priority() int {
	return 20
}
