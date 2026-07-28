package settings

import "Paarthurnax/internal/domain/translation"

type Settings struct {
	SourceLocale       translation.Locale
	DestinationLocales []translation.Locale
}
