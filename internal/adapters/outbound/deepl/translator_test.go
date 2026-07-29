package deepl

import (
	"Paarthurnax/internal/domain/translation"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var protectedVariablePattern = regexp.MustCompile(`<span translate="no" id="[0-9a-f-]{36}">([^<]+)</span>`)

type fakeClient struct {
	text              string
	sourceLocale      string
	destinationLocale string
	response          func(string) string
	err               error
}

func (c *fakeClient) Translate(text string, sourceLocale string, destinationLocale string) (string, error) {
	c.text = text
	c.sourceLocale = sourceLocale
	c.destinationLocale = destinationLocale
	if c.err != nil {
		return "", c.err
	}
	return c.response(text), nil
}

func TestTranslatorTranslate(t *testing.T) {
	t.Parallel()

	client := &fakeClient{
		response: func(text string) string {
			return strings.Replace(text, "Hello", "Bonjour", 1)
		},
	}
	translator := NewTranslator(client)
	segment := *translation.NewSegment("greeting", "Hello %{name}!")

	got, err := translator.Translate(segment, translation.English, translation.French)

	if err != nil {
		t.Fatalf("Translate() returned unexpected error: %v", err)
	}
	if client.sourceLocale != "en" || client.destinationLocale != "fr" {
		t.Fatalf("Translate() locales = (%q, %q), want (%q, %q)", client.sourceLocale, client.destinationLocale, "en", "fr")
	}
	if !protectedVariablePattern.MatchString(client.text) {
		t.Fatalf("Translate() sent unprotected variables to client: %q", client.text)
	}

	want := translation.NewSegment("greeting", "Bonjour %{name}!")
	if !reflect.DeepEqual(got, *want) {
		t.Fatalf("Translate() = %#v, want %#v", got, *want)
	}
}

func TestTranslatorTranslateClientError(t *testing.T) {
	t.Parallel()

	clientErr := errors.New("DeepL unavailable")
	translator := NewTranslator(&fakeClient{err: clientErr})
	segment := *translation.NewSegment("greeting", "Hello %{name}!")

	got, err := translator.Translate(segment, translation.English, translation.French)

	if !errors.Is(err, clientErr) {
		t.Fatalf("Translate() error = %v, want %v", err, clientErr)
	}
	if !reflect.DeepEqual(got, segment) {
		t.Fatalf("Translate() = %#v, want original segment %#v", got, segment)
	}
}

func TestBuildTranslationText(t *testing.T) {
	t.Parallel()

	segment := translation.NewSegment("greeting", "Hello %{name}, you have %{count} messages for %{name}.")

	text, reverseHash := buildTranslationText(*segment)

	matches := protectedVariablePattern.FindAllStringSubmatch(text, -1)
	if len(matches) != 3 {
		t.Fatalf("buildTranslationText() generated %d protected variables, want 3: %q", len(matches), text)
	}

	gotNames := []string{matches[0][1], matches[1][1], matches[2][1]}
	wantNames := []string{"name", "count", "name"}
	for i := range wantNames {
		if gotNames[i] != wantNames[i] {
			t.Fatalf("protected variable %d = %q, want %q", i, gotNames[i], wantNames[i])
		}
	}

	if len(reverseHash) != 3 {
		t.Fatalf("reverse hash has %d entries, want 3", len(reverseHash))
	}
	if got := reverseEncodedText(text, reverseHash); got != segment.Value {
		t.Fatalf("reverseEncodedText() = %q, want %q", got, segment.Value)
	}
}

func TestBuildTranslationTextPlainText(t *testing.T) {
	t.Parallel()

	segment := translation.NewSegment("greeting", "Hello world")

	text, reverseHash := buildTranslationText(*segment)

	if text != "Hello world" {
		t.Fatalf("buildTranslationText() = %q, want %q", text, "Hello world")
	}
	if len(reverseHash) != 0 {
		t.Fatalf("reverse hash has %d entries, want 0", len(reverseHash))
	}
}

func TestReverseEncodedText(t *testing.T) {
	t.Parallel()

	const encodedName = `<span translate="no" id="name-id">name</span>`
	const encodedCount = `<span translate="no" id="count-id">count</span>`

	got := reverseEncodedText(
		"Bonjour "+encodedName+", vous avez "+encodedCount+" messages. "+encodedName,
		map[string]string{encodedName: "name", encodedCount: "count"},
	)
	want := "Bonjour %{name}, vous avez %{count} messages. %{name}"
	if got != want {
		t.Fatalf("reverseEncodedText() = %q, want %q", got, want)
	}
}
