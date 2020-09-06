package svg

import "regexp"

type FormatWeekdayHeader struct {
	Attribute
}

func (f FormatWeekdayHeader) Apply(text CalendarText) {
	if text.WeekdayHeader == 0 {
		return
	}

	format := f.Attr().(string)
	switch format {
	case "1":
		text.Content = text.Content[:1]
		break
	case "2":
	case "2p":
		// todo: lang
		text.Content = text.Content[:2]
		if format == "2p" {
			text.Content += "."
		}
		break
	}
}

func (f FormatWeekdayHeader) Matches(subject string) bool {
	reg := regexp.MustCompile("^dt([129p]{1,2})")
	return reg.MatchString(subject)
}

func (f FormatWeekdayHeader) New(subject string) Annotation {
	reg := regexp.MustCompile("^dt([129p]{1,2})")
	val := reg.FindString(subject)
	return FormatWeekdayHeader{Attribute{val}}
}

func (f FormatWeekdayHeader) Id() string {
	return "formatWeekdayHeader"
}

func (f FormatWeekdayHeader) Priority() int {
	return 20
}

