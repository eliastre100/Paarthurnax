package sync

type ChangeType uint8

const (
	ChangeTypeInsert ChangeType = iota
	ChangeTypeDelete
	ChangeTypeUpdate
	ChangeTypeCreateFile
	ChangeTypeDeleteFile
)

type Change struct {
	Type     ChangeType
	Document *Document
	Segment  *string
}

type DocumentChanges struct {
	Document *Document
	Changes  []*Change
}
