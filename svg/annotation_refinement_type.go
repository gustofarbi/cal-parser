package svg

import "regexp"

type RefinementType struct {
	Attribute
}

func (r RefinementType) Apply(text *CalendarText) {}

func (r RefinementType) Matches(subject string) bool {
	reg := regexp.MustCompile("ref=([glrs])")
	return reg.MatchString(subject)
}

func (r RefinementType) New(subject string) Annotation {
	reg := regexp.MustCompile("ref=([glrs])")
	raw := reg.FindStringSubmatch(subject)
	return RefinementType{Attribute{raw[1]}}
}

func (r RefinementType) Id() string {
	return "refinementType"
}

func (r RefinementType) Priority() int {
	return 0
}

