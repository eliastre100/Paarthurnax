package translation

import (
	"regexp"
	"strings"
)

type Segment struct {
	Key    string
	Value  string
	Plural bool
	Parts  []SegmentPart
}

type SegmentPart interface {
}

type VariablePart struct {
	Name string
}

type TextPart struct {
	Text string
}

var variablePattern = regexp.MustCompile(`%\{\s*([a-zA-Z_ -]+)\s*}`)
var pluralCategories = map[string]struct{}{
	"zero":  {},
	"one":   {},
	"two":   {},
	"few":   {},
	"many":  {},
	"other": {},
}

func NewSegment(key string, value string) *Segment {
	parts := parseSegmentParts(value)
	segment := &Segment{
		Key:   key,
		Value: value,
		Parts: parts,
	}

	segment.Plural = isPluralKey(key) && segment.HasVariable("count")

	return segment
}

func isPluralKey(key string) bool {
	lastSeparator := strings.LastIndex(key, ".")
	if lastSeparator == -1 || lastSeparator == len(key)-1 {
		return false
	}

	_, ok := pluralCategories[key[lastSeparator+1:]]
	return ok
}

func (segment *Segment) HasVariable(name string) bool {
	for _, part := range segment.Parts {
		variable, ok := part.(VariablePart)
		if ok && variable.Name == name {
			return true
		}
	}

	return false
}

func parseSegmentParts(value string) []SegmentPart {
	matches := variablePattern.FindAllStringSubmatchIndex(value, -1)
	if len(matches) == 0 {
		return []SegmentPart{TextPart{Text: value}}
	}

	// Each placeholder can contribute one VariablePart and one following TextPart
	// after it, add the first potential TextPart, so 2*n+1 is the upper bound for n matches.
	parts := make([]SegmentPart, 0, len(matches)*2+1)
	last := 0
	for _, match := range matches {
		start, end := match[0], match[1]
		if start > last {
			parts = append(parts, TextPart{Text: value[last:start]})
		}

		nameStart, nameEnd := match[2], match[3]
		parts = append(parts, VariablePart{Name: strings.TrimSpace(value[nameStart:nameEnd])})
		last = end
	}

	if last < len(value) {
		parts = append(parts, TextPart{Text: value[last:]})
	}

	return parts
}
