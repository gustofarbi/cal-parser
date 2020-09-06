package svg

type Calendar struct {
	texts                []CalendarText
	weekdayHeadingsTable []CalendarText
	weekdayHeadingsLine  []CalendarText
	Receiver             chan interface{}
	Context              Context
	RenderPrevNext       bool
}

func (c Calendar) StartReceiver() {
	for {
		select {
		//case object := <- c.Receiver:
		//	switch x := object.(type) {
		//	todo
		//	}
		}
	}
}

// todo: font-style
type CalendarText struct {
	Position
	Annotations AnnotationCollection
	Content string
	FontSize float64
	FontFamily string
	FontColor string
	CalendarType string
	IsYear bool
	IsMonth bool
	CalendarWeek int
	WeekdayHeader int
	WeekdayPosition int
	IsWeekend bool
	Language string
	IsCurrentMonth bool
	FormatWeekdayPosition string
	FormatMonth string
	FormatYear string
	HasRefinement bool
}
