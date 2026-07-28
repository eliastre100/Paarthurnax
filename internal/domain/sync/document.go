package sync

import "errors"

type Digest string

type Document struct {
	Name    string
	Digests map[string]Digest
}

var ErrSegmentAlreadyExists = errors.New("segment already exists")

func NewDocument(name string) *Document {
	return &Document{
		Name:    name,
		Digests: make(map[string]Digest),
	}
}

func (d *Document) AddSegmentDigest(segment string, digest Digest) error {
	if _, ok := d.Digests[segment]; ok {
		return ErrSegmentAlreadyExists
	}
	d.Digests[segment] = digest
	return nil
}

func (d *Document) GetSegmentDigest(segment string) Digest {
	return d.Digests[segment]
}
