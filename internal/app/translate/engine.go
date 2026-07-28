package translate

import "Paarthurnax/internal/domain/translation"

type SegmentValue struct {
	Name  string
	Value string
}

type Engine interface {
	Translate(segment translation.Segment, srcLocale translation.Locale, dstLocale translation.Locale, values ...SegmentValue) (translation.Segment, error)
}
