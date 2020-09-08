package svg

import (
	"fmt"
	"sync"
	"time"
)

type Calendar struct {
	texts                     []CalendarText
	weekdayHeadingsTable      []CalendarText
	weekdayHeadingsLine       []CalendarText
	positionTableCurrentMonth []CalendarText
	positionTableAnotherMonth []CalendarText
	positionsLineWeekend      []CalendarText
	positionsLineWeekday      []CalendarText
	positionsLineDefault      []CalendarText
	years                     []CalendarText
	months                    []CalendarText
	calendarWeeksTable        []CalendarText
	calendarWeeksLine         []CalendarText

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
		weekdayHeadingsTable:      make([]CalendarText, 0),
		weekdayHeadingsLine:       make([]CalendarText, 0),
		positionTableCurrentMonth: make([]CalendarText, 0),
		positionTableAnotherMonth: make([]CalendarText, 0),
		positionsLineWeekend:      make([]CalendarText, 0),
		positionsLineWeekday:      make([]CalendarText, 0),
		positionsLineDefault:      make([]CalendarText, 0),
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
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			fmt.Println("still running")
		case item := <-c.Receiver:
			switch x := item.(type) {
			case CalendarText:
				fmt.Println("received text")
				c.SaveText(x)
			case AnnotationObject:
				fmt.Println("received annotation object")
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
			c.weekdayHeadingsTable = append(c.weekdayHeadingsTable, x)
		} else if x.CalendarType == "line" {
			c.weekdayHeadingsLine = append(c.weekdayHeadingsLine, x)
		}
	case x.WeekdayPosition > 0:
		if x.CalendarType == "table" {
			if x.IsCurrentMonth { // todo: switch these
				c.positionTableAnotherMonth = append(c.positionTableAnotherMonth, x)
			} else {
				c.positionTableCurrentMonth = append(c.positionTableCurrentMonth, x)
			}
		} else if x.CalendarType == "line" {
			if x.IsWeekend {
				c.positionsLineWeekend = append(c.positionsLineWeekend, x)
			} else if x.IsWeekday {
				c.positionsLineWeekday = append(c.positionsLineWeekday, x)
			} else {
				c.positionsLineDefault = append(c.positionsLineDefault, x)
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
}
