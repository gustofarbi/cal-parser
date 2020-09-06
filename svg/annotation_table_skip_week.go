package svg

import (
	"regexp"
	"strconv"
)

type SkipWeek struct {
	Attribute
}

func (s SkipWeek) Apply(text CalendarText) {}

func (s SkipWeek) Matches(subject string) bool {
	reg := regexp.MustCompile("^w([56])")
	return reg.MatchString(subject)
}

func (s SkipWeek) New(subject string) Annotation {
	reg := regexp.MustCompile("^w([56])")
	val := reg.FindString(subject)
	attr, _ := strconv.Atoi(val)
	return SkipWeek{Attribute{attr}}
}

func (s SkipWeek) Id() string {
	return "skipWeek"
}

func (s SkipWeek) Priority() int {
	return 0
}
