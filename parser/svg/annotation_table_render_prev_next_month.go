package svg

import (
	"strconv"
	"strings"
)

type RenderPrevNextMonth struct {
	Attribute
}

func (r RenderPrevNextMonth) Apply(text *CalendarText) {}

func (r RenderPrevNextMonth) Matches(subject string) bool {
	return strings.Contains(subject, "mpn=")
}

func (r RenderPrevNextMonth) New(subject string) Annotation {
	val, _ := strconv.ParseBool(subject[4:])
	return RenderPrevNextMonth{Attribute{val}}
}

func (r RenderPrevNextMonth) Id() string {
	return "renderPrevNextMonth"
}

func (r RenderPrevNextMonth) Priority() int {
	return 0
}
