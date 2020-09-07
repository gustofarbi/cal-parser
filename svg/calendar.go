package svg

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

	skipWeeks       map[int]string
	skipDays        map[int]string
	weekdayElements map[int]string
	weekendElements map[int]string

	Receiver       chan interface{}
	RenderPrevNext bool
	svgContent     string
}

func (c Calendar) StartReceiver() {
	for {
		select {
		case item := <-c.Receiver:
			switch x := item.(type) {
			case CalendarText:
				c.SaveText(x)
			case AnnotationObject:
				x.
			}
		}
	}
}

func (c Calendar) SaveText(x CalendarText) {
	switch {
	case x.WeekdayHeader > 0:
		if x.CalendarType == "table" {
			c.weekdayHeadingsTable = append(c.weekdayHeadingsTable, x)
		} else if x.CalendarType == "line" {
			c.weekdayHeadingsLine = append(c.weekdayHeadingsLine, x)
		}
	case x.WeekdayPosition > 0:
		if x.CalendarType == "table" {
			if x.IsCurrentMonth {
				c.positionTableCurrentMonth = append(c.positionTableCurrentMonth, x)
			} else {
				c.positionTableAnotherMonth = append(c.positionTableAnotherMonth, x)
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
