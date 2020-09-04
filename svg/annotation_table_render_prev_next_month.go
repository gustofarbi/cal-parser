package svg

import (
	"strconv"
	"strings"
)

type RenderPrevNextMonth struct {
	attr bool
}

func (r RenderPrevNextMonth) Apply(text CalendarText) {}

func (r RenderPrevNextMonth) Matches(subject string) bool {
	return strings.Contains(subject, "npm=")
}

func (r RenderPrevNextMonth) New(subject string) Annotation {
	val, _ := strconv.ParseBool(subject[4:])
	return RenderPrevNextMonth{val}
}

func (r RenderPrevNextMonth) Id() string {
	return "renderPrevNextMonth"
}

func (r RenderPrevNextMonth) Priority() int {
	return 0
}
