package svg

import "regexp"

var months = map[int]string{
	1:  "Januar",
	2:  "Februar",
	3:  "März",
	4:  "April",
	5:  "Mai",
	6:  "Juni",
	7:  "Juli",
	8:  "August",
	9:  "September",
	10: "Oktober",
	11: "November",
	12: "Dezember",
}

type FormatMonthText struct {
	Attribute
}

func (f FormatMonthText) Apply(text *CalendarText) {
	if !text.IsMonth {
		return
	}
	text.Content = months[text.CurrentMonth]
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
