package reporter

import (
	"bytes"
	"strings"
	"testing"

	"charm.land/log/v2"
)

func TestTextTracksDocumentProgress(t *testing.T) {
	output := captureLogOutput(t)
	reporter := NewTextReporter()

	reporter.UpdateDocumentChanges("messages", 2)
	reporter.DoneUpdatingSegment("messages", "welcome")
	reporter.DoneDeletingSegment("messages", "obsolete")

	mustContainLog(t, output, "Translated welcome (1 / 2)")
	mustContainLog(t, output, "Deleted obsolete (2 / 2)")
}

func TestTextUpdatesExistingDocumentProgressTotal(t *testing.T) {
	output := captureLogOutput(t)
	reporter := NewTextReporter()

	reporter.UpdateDocumentChanges("messages", 1)
	reporter.UpdateDocumentChanges("messages", 3)
	reporter.DoneUpdatingSegment("messages", "welcome")

	mustContainLog(t, output, "Translated welcome (1 / 3)")
}

func mustContainLog(t *testing.T, output *bytes.Buffer, want string) {
	t.Helper()

	if !strings.Contains(output.String(), want) {
		t.Fatalf("log output = %q, want it to contain %q", output.String(), want)
	}
}

func captureLogOutput(t *testing.T) *bytes.Buffer {
	t.Helper()

	previousLogger := log.Default()
	output := &bytes.Buffer{}
	log.SetDefault(log.New(output))
	t.Cleanup(func() {
		log.SetDefault(previousLogger)
	})

	return output
}
