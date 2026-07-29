package deepl

import (
	"Paarthurnax/internal/app/translate"
	"Paarthurnax/internal/domain/translation"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Translator struct {
	client translationClient
}

type translationClient interface {
	Translate(text string, sourceLocale string, destinationLocale string) (string, error)
}

func NewTranslator(client translationClient) *Translator {
	return &Translator{client: client}
}

func (t *Translator) Translate(segment translation.Segment, srcLocale translation.Locale, dstLocale translation.Locale, values ...translate.SegmentValue) (translation.Segment, error) {
	txtSegment, reverseHash := buildTranslationText(segment, values...)
	result, err := t.client.Translate(txtSegment, srcLocale.String(), dstLocale.String())
	if err != nil {
		return segment, err
	}
	translatedSegment := translation.NewSegment(segment.Key, reverseEncodedText(result, reverseHash))

	return *translatedSegment, nil
}

func buildTranslationText(segment translation.Segment, values ...translate.SegmentValue) (string, map[string]string) {
	txt := ""
	reverseHash := make(map[string]string)

	for _, part := range segment.Parts {
		switch part.(type) {
		case translation.TextPart:
			txt += part.(translation.TextPart).Text
		case translation.VariablePart:
			id := uuid.New().String()
			encoded := fmt.Sprintf("<span translate=\"no\" id=\"%s\">%s</span>", id, part.(translation.VariablePart).Name)
			reverseHash[encoded] = part.(translation.VariablePart).Name
			txt += encoded
		}
	}

	return txt, reverseHash
}

func reverseEncodedText(txt string, reverseHash map[string]string) string {
	for encoded, variable := range reverseHash {
		txt = strings.ReplaceAll(txt, encoded, "%{"+variable+"}")
	}
	return txt
}
