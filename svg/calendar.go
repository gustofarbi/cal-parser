package svg

type Calendar struct {
	texts []CalendarText
	weekdayHeadingsTable []CalendarText
	weekdayHeadingsLine []CalendarText

	Context Context
}

// todo: font-style
type CalendarText struct {
	Position
	Content string
	FontSize float64
	FontFamily string
	FontColor string
	CalendarType string
}
