package reporter

import (
	"fmt"
	"sync"

	"charm.land/log/v2"
)

type DocumentProgress struct {
	Done  uint
	Total uint
}

type Text struct {
	mutex    sync.Mutex
	progress map[string]DocumentProgress
}

func NewTextReporter() *Text {
	return &Text{
		progress: make(map[string]DocumentProgress),
	}
}

func (r *Text) StartHandlingDocument(document string) {
	log.Infof("Handling %s", document)
}

func (r *Text) StartHandlingChanges(documentCount uint) {
}

func (r *Text) UpdateDocumentChanges(document string, count uint) {
	r.updateDocumentProgressTotal(document, count)
}

func (r *Text) DoneUpdatingSegment(document string, name string) {
	progress := r.incrementDocumentProgress(document)
	log.Info(fmt.Sprintf("Translated %s (%d / %d)", name, progress.Done, progress.Total), "document", document)
}

func (r *Text) DoneUpdatingDocument(document string) {
}

func (r *Text) DoneDeletingDocument(document string) {
	log.Infof("Cleaning up %s", document)
}

func (r *Text) DoneDeletingSegment(document string, name string) {
	progress := r.incrementDocumentProgress(document)
	log.Info(fmt.Sprintf("Deleted %s (%d / %d)", name, progress.Done, progress.Total), "document", document)
}

func (r *Text) updateDocumentProgressTotal(document string, total uint) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	progress, ok := r.progress[document]
	if !ok {
		progress = DocumentProgress{Total: total}
	}

	progress.Total = total
	r.progress[document] = progress
}

func (r *Text) incrementDocumentProgress(document string) DocumentProgress {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	progress, ok := r.progress[document]
	if !ok {
		progress = DocumentProgress{}
	}

	progress.Done = progress.Done + 1
	r.progress[document] = progress
	return progress
}
