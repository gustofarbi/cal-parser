package svg

var languages = []string{
	"de",
	"fr",
	"nl",
	"en",
	"it",
	"es",
}

type Language struct {
	Lang string
}

func (l Language) Apply(text CalendarText) {
	text.Language = l.Lang
}

func (l Language) Matches(subject string) bool {
	for _, lang := range languages {
		if subject == lang {
			return true
		}
	}
	return false
}

func (l Language) New(subject string) Annotation {
	return Language{subject}
}

func (l Language) Id() string {
	return "language"
}

func (l Language) Priority() int {
	return 0
}
