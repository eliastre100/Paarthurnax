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
	translatedText, err := reverseEncodedText(result, reverseHash)
	if err != nil {
		return segment, fmt.Errorf("unable to restore protected variables: %w", err)
	}

	translatedSegment := translation.NewSegment(segment.Key, translatedText)
	if err := checkVariableEquity(segment, *translatedSegment); err != nil {
		return segment, fmt.Errorf("translated segment variables do not match the source: %w", err)
	}

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

func reverseEncodedText(txt string, reverseHash map[string]string) (string, error) {
	for encoded, variable := range reverseHash {
		if count := strings.Count(txt, encoded); count != 1 {
			return "", fmt.Errorf("protected variable %q occurs %d times, want 1", variable, count)
		}
		txt = strings.ReplaceAll(txt, encoded, "%{"+variable+"}")
	}
	return txt, nil
}

func checkVariableEquity(source translation.Segment, translated translation.Segment) error {
	sourceVariables := variableCounts(source)
	translatedVariables := variableCounts(translated)

	for variable, count := range sourceVariables {
		if translatedVariables[variable] != count {
			return fmt.Errorf("variable %q occurs %d times in translation, want %d", variable, translatedVariables[variable], count)
		}
	}
	for variable, count := range translatedVariables {
		if sourceVariables[variable] != count {
			return fmt.Errorf("unexpected variable %q occurs %d times in translation", variable, count)
		}
	}

	return nil
}

func variableCounts(segment translation.Segment) map[string]int {
	counts := make(map[string]int)
	for _, part := range segment.Parts {
		if variable, ok := part.(translation.VariablePart); ok {
			counts[variable.Name]++
		}
	}
	return counts
}
