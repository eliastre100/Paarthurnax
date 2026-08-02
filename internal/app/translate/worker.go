package translate

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"
	"errors"
	"fmt"
)

type Worker struct {
	settings      *settings.Settings
	project       *translation.Project
	engine        Engine
	documentStore DocumentStore
	reporter      Reporter
}

func NewWorker(settings *settings.Settings, project *translation.Project, engine Engine, documentStore DocumentStore, reporter Reporter) *Worker {
	return &Worker{
		settings:      settings,
		project:       project,
		engine:        engine,
		documentStore: documentStore,
		reporter:      reporter,
	}
}

func (w *Worker) Handle(documentChanges sync.DocumentChanges) error {
	// Notify start working on document
	w.reporter.StartHandlingDocument(documentChanges.Document.Name)
	w.reporter.UpdateDocumentChanges(documentChanges.Document.Name, uint(len(documentChanges.Changes)))
	for _, change := range documentChanges.Changes {
		switch change.Type {
		case sync.ChangeTypeDeleteFile:
			return w.RemoveDocument(documentChanges.Document.Name)
		case sync.ChangeTypeCreateFile:
			if err := w.createDocument(documentChanges.Document.Name); err != nil {
				return err
			}
		case sync.ChangeTypeDelete:
			if err := w.RemoveSegment(documentChanges.Document.Name, *change.Segment); err != nil {
				return err
			}
		case sync.ChangeTypeUpdate, sync.ChangeTypeInsert:
			if err := w.UpdateSegment(documentChanges.Document.Name, *change.Segment); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown change type %d", change.Type)
		}
	}

	return w.persistChanges(documentChanges.Document.Name)
}

func (w *Worker) createDocument(document string) error {
	sourceDocument := w.project.Documents[document]
	catalog := sourceDocument.Catalog(w.settings.SourceLocale)

	w.reporter.UpdateDocumentChanges(document, uint(len(catalog.Segments)))
	for key := range catalog.Segments {
		if err := w.UpdateSegment(document, key); err != nil {
			return fmt.Errorf("failed to update segment %s in document %s: %w", key, document, err)
		}
	}
	return nil
}

func (w *Worker) RemoveSegment(documentName string, key string) error {
	// TODO: handle plural key deletion
	for _, locale := range w.settings.DestinationLocales {
		document, err := w.project.GetDocumentIn(locale, documentName)
		if err != nil {
			return fmt.Errorf("failed to remove segment %s in document %s: %w", key, documentName, err)
		}
		document.RemoveSegment(key)
	}

	w.reporter.DoneDeletingSegment(documentName, key)
	return nil
}

func (w *Worker) RemoveDocument(name string) error {
	document := translation.NewDocument(name)
	if err := document.AddCatalog(translation.NewCatalog(w.settings.SourceLocale)); err != nil {
		return err
	}

	for _, locale := range w.settings.DestinationLocales {
		localeDocumentName := document.DocumentNameForLocale(locale)
		if err := w.documentStore.Delete(localeDocumentName); err != nil {
			return fmt.Errorf("failed to delete document %s for locale %s: %w", localeDocumentName, locale, err)
		}
	}

	w.reporter.DoneDeletingDocument(name)
	return nil
}

func (w *Worker) UpdateSegment(documentName string, key string) error {
	document, ok := w.project.Documents[documentName]
	if !ok {
		return fmt.Errorf("document %s not found", documentName)
	}
	src, err := document.GetSegment(w.settings.SourceLocale, key)
	if err != nil {
		return fmt.Errorf("segment %s not found in document %s: %w", key, documentName, err)
	}

	for _, locale := range w.settings.DestinationLocales {
		result, err := w.engine.Translate(src, w.settings.SourceLocale, locale) // TODO: inject count value for plurals & do it for every target key
		if err != nil {
			return fmt.Errorf("failed to translate segment %s in document %s from %s to %s: %w", key, documentName, w.settings.SourceLocale, locale, err)
		}
		targetDocument, err := w.project.GetDocumentIn(locale, documentName)
		if err != nil {
			if errors.Is(err, translation.ErrDocumentNotFound) {
				destName := document.DocumentNameForLocale(locale)
				createdDocument := translation.NewDocument(destName)
				targetDocument = createdDocument
				w.project.AddDocument(targetDocument)
			} else {
				return fmt.Errorf("failed to locate destination document for segment %s in document %s and locale %s: %w", key, documentName, locale, err)
			}
		}
		if err := targetDocument.SetSegment(locale, result); err != nil {
			return fmt.Errorf("failed to set segment %s in document %s: %w", key, targetDocument.Name, err)
		}
	}

	w.reporter.DoneUpdatingSegment(documentName, key)
	return nil
}

func (w *Worker) persistChanges(name string) error {
	for _, locale := range w.settings.DestinationLocales {
		document, err := w.project.GetDocumentIn(locale, name)
		if err != nil {
			return fmt.Errorf("failed to locate document %s for locale %s before persistence: %w", name, locale, err)
		}
		if err := w.documentStore.Store(document); err != nil {
			return fmt.Errorf("failed to store document %s for locale %s: %w", document.Name, locale, err)
		}
	}
	w.reporter.DoneUpdatingDocument(name)
	return nil
}
