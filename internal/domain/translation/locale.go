package translation

type Locale string

const (
	English  Locale = "en"
	French   Locale = "fr"
	Japanese Locale = "ja"

	// TODO: Add the remaining locales
)

var Locales = []Locale{English, French}
