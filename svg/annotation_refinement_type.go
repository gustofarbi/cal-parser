package svg

import "regexp"

type RefinementType struct {
	Type string
}

func (r RefinementType) Apply(text CalendarText) {}

func (r RefinementType) Matches(subject string) bool {
	reg := regexp.MustCompile("ref=([glrs])")
	return reg.MatchString(subject)
}

func (r RefinementType) New(subject string) Annotation {
	return RefinementType{subject}
}

func (r RefinementType) Id() string {
	return "refinementType"
}

func (r RefinementType) Priority() int {
	return 0
}

