package svg

import (
	"sync"
	"time"
)

type Calendar struct {
	texts                []CalendarText
	weekdayHeadingsTable map[int]CalendarText
	weekdayHeadingsLine  map[int]CalendarText

	positionTableCurrentMonth map[int]CalendarText
	positionTableAnotherMonth map[int]CalendarText
	positionsLineWeekend      map[int]CalendarText
	positionsLineWeekday      map[int]CalendarText
	positionsLineDefault      map[int]CalendarText

	calendarWeeksTable []CalendarText
	calendarWeeksLine  []CalendarText

	years  []CalendarText
	months []CalendarText

	skipWeeks        map[int][]string
	skipDays         map[int][]string
	weekdayElements  map[int][]string
	weekendElements  map[int][]string
	renderMonthsOnly map[int][]string

	Receiver       chan interface{}
	ReceiverWg     *sync.WaitGroup
	RenderPrevNext bool
	svgContent     string
}

func NewCalendar() Calendar {
	return Calendar{
		texts:                     make([]CalendarText, 0),
		weekdayHeadingsTable:      make(map[int]CalendarText),
		weekdayHeadingsLine:       make(map[int]CalendarText),
		positionTableCurrentMonth: make(map[int]CalendarText),
		positionTableAnotherMonth: make(map[int]CalendarText),
		positionsLineWeekend:      make(map[int]CalendarText),
		positionsLineWeekday:      make(map[int]CalendarText),
		positionsLineDefault:      make(map[int]CalendarText),
		years:                     make([]CalendarText, 0),
		months:                    make([]CalendarText, 0),
		calendarWeeksTable:        make([]CalendarText, 0),
		calendarWeeksLine:         make([]CalendarText, 0),

		skipWeeks:        make(map[int][]string),
		skipDays:         make(map[int][]string),
		weekdayElements:  make(map[int][]string),
		weekendElements:  make(map[int][]string),
		renderMonthsOnly: make(map[int][]string),

		Receiver:   make(chan interface{}),
		ReceiverWg: &sync.WaitGroup{},

		RenderPrevNext: false,
		svgContent:     "",
	}
}

func (c *Calendar) StartReceiver() {
	ticker := time.Tick(100 * time.Millisecond)
	counter := 0
loop:
	for {
		select {
		case <-ticker:
			counter++
			if counter > 10 {
				break loop
			}
		case item := <-c.Receiver:
			counter = 0
			switch x := item.(type) {
			case CalendarText:
				c.SaveText(x)
			case AnnotationObject:
				switch o := x.Annotation.(type) {
				case RenderPrevNextMonth:
					if !c.RenderPrevNext {
						c.RenderPrevNext = x.Attr().(bool)
					}
				case SkipWeek:
					c.skipWeeks[o.Attr().(int)] = append(c.skipWeeks[o.Attr().(int)], x.RawSvg)
				case LineSkipDay:
					c.skipDays[o.Attr().(int)] = append(c.skipDays[o.Attr().(int)], x.RawSvg)
				case LineWeekdayElement:
					c.weekdayElements[o.Attr().(int)] = append(c.weekdayElements[o.Attr().(int)], x.RawSvg)
				case LineWeekendElement:
					c.weekendElements[o.Attr().(int)] = append(c.weekendElements[o.Attr().(int)], x.RawSvg)
				}
			}
		}
	}
}

func (c *Calendar) SaveText(x CalendarText) {
	switch {
	case x.WeekdayHeader > 0:
		if x.CalendarType == "table" {
			c.weekdayHeadingsTable[x.WeekdayHeader] = x
		} else if x.CalendarType == "line" {
			c.weekdayHeadingsLine[x.WeekdayHeader] = x
		}
	case x.WeekdayPosition > 0:
		if x.CalendarType == "table" {
			if x.IsCurrentMonth { // todo: switch these
				c.positionTableAnotherMonth[x.WeekdayPosition] = x
			} else {
				c.positionTableCurrentMonth[x.WeekdayPosition] = x
			}
		} else if x.CalendarType == "line" {
			if x.IsWeekend {
				c.positionsLineWeekend[x.WeekdayPosition] = x
			} else if x.IsWeekday {
				c.positionsLineWeekday[x.WeekdayPosition] = x
			} else {
				c.positionsLineDefault[x.WeekdayPosition] = x
			}
		}
	case x.IsYear:
		c.years = append(c.years, x)
	case x.IsMonth:
		c.months = append(c.months, x)
	case x.CalendarWeek > 0:
		if x.CalendarType == "table" {
			c.calendarWeeksTable = append(c.calendarWeeksTable, x)
		} else {
			c.calendarWeeksLine = append(c.calendarWeeksLine, x)
		}
	default:
		c.texts = append(c.texts, x)
	}
}

// todo: font-style
type CalendarText struct {
	Position
	Annotations           AnnotationCollection
	Content               string
	FontSize              float64
	FontFamily            string
	FontColor             string
	CalendarType          string
	IsYear                bool
	IsMonth               bool
	CalendarWeek          int
	WeekdayHeader         int
	WeekdayPosition       int
	IsWeekend             bool
	IsWeekday             bool
	Language              string
	IsCurrentMonth        bool
	FormatWeekdayPosition string
	FormatMonth           string
	FormatYear            string
	HasRefinement         bool
	CurrentMonth          int
	CurrentYear           int
}
