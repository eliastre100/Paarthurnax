package translation

import (
	"errors"

	"github.com/google/uuid"
)

type Catalog struct {
	ID       string
	Locale   Locale
	Segments map[string]Segment
}

var ErrSegmentNotFound = errors.New("segment not found")

func NewCatalog(locale Locale) *Catalog {
	return &Catalog{
		ID:       uuid.New().String(),
		Locale:   locale,
		Segments: make(map[string]Segment),
	}
}

func (c *Catalog) AddSegment(segment Segment) {
	c.Segments[segment.Key] = segment
}

func (c *Catalog) RemoveSegment(key string) {
	delete(c.Segments, key)
}

func (c *Catalog) UpdateSegment(key string, segment Segment) error {
	_, ok := c.Segments[key]
	if !ok {
		return ErrSegmentNotFound
	}
	c.Segments[key] = segment
	return nil
}

func (c *Catalog) GetSegment(key string) (Segment, error) {
	segment, ok := c.Segments[key]
	if !ok {
		return Segment{}, ErrSegmentNotFound
	}
	return segment, nil
}

func (c *Catalog) isEmpty() bool {
	return len(c.Segments) == 0
}
