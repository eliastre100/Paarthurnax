package translation

import (
	"slices"
	"strconv"
	"testing"
)

func TestLocalePluralCodex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		source Locale
		target Locale
		want   map[string][]Plural
	}{
		{
			name:   "French to Romanian",
			source: French,
			target: Romanian,
			want: map[string][]Plural{
				"one":   {{Key: "one", Hint: "1"}},
				"other": {{Key: "few", Hint: "2"}, {Key: "other", Hint: "20"}},
				"zero":  {{Key: "zero", Hint: "0"}},
			},
		},
		{
			name:   "english to Arabic maps zero to the reserved zero key",
			source: English,
			target: Arabic,
			want: map[string][]Plural{
				"one":  {{Key: "one", Hint: "1"}},
				"zero": {{Key: "zero", Hint: "0"}},
				"other": {
					{Key: "two", Hint: "2"},
					{Key: "few", Hint: "3"},
					{Key: "many", Hint: "11"},
					{Key: "other", Hint: "100"},
				},
			},
		},
		{
			name:   "english to Japanese",
			source: English,
			target: Japanese,
			want: map[string][]Plural{
				"one":  {{Key: "other", Hint: "1"}},
				"zero": {{Key: "zero", Hint: "0"}},
			},
		},
		{
			name:   "Romanian to French",
			source: Romanian,
			target: French,
			want: map[string][]Plural{
				"one":  {{Key: "one", Hint: "1"}},
				"few":  {{Key: "other", Hint: "2"}},
				"zero": {{Key: "zero", Hint: "0"}},
			},
		},
		{
			name:   "Arabic to Arabic preserves all categories",
			source: Arabic,
			target: Arabic,
			want: map[string][]Plural{
				"zero":  {{Key: "zero", Hint: "0"}},
				"one":   {{Key: "one", Hint: "1"}},
				"two":   {{Key: "two", Hint: "2"}},
				"few":   {{Key: "few", Hint: "3"}},
				"many":  {{Key: "many", Hint: "11"}},
				"other": {{Key: "other", Hint: "100"}},
			},
		},
		{
			name:   "unknown source locale",
			source: "unknown",
			target: English,
			want:   map[string][]Plural{},
		},
		{
			name:   "unknown target locale",
			source: English,
			target: "unknown",
			want:   map[string][]Plural{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.source.PluralCodex(tt.target); !equalPluralCodex(got, tt.want) {
				t.Errorf("PluralCodex(%q) = %#v, want %#v", tt.target, got, tt.want)
			}
		})
	}
}

func TestLocalePluralCodexSupportsAllDeclaredLocalePairs(t *testing.T) {
	t.Parallel()

	for _, source := range Locales {
		for _, target := range Locales {
			sourceDefinition, sourceSupported := plurals[source.String()]
			targetDefinition, targetSupported := plurals[target.String()]
			codex := source.PluralCodex(target)

			if !sourceSupported || !targetSupported {
				if len(codex) != 0 {
					t.Errorf("PluralCodex(%q, %q) = %#v, want an empty codex for an unsupported locale", source, target, codex)
				}
				continue
			}

			want := pluralCodex(sourceDefinition, targetDefinition)
			if !equalPluralCodex(codex, want) {
				t.Errorf("PluralCodex(%q, %q) = %#v, want %#v", source, target, codex, want)
			}
		}
	}
}

func pluralCodex(source, target definition) map[string][]Plural {
	codex := make(map[string][]Plural)
	for _, hint := range target.Hints {
		if hint.Key == "zero" {
			continue
		}

		sourceKey := string(source.Evaluator(hint.Hint))
		targetPlural := Plural{Key: hint.Key, Hint: strconv.Itoa(hint.Hint)}
		codex[sourceKey] = append(codex[sourceKey], targetPlural)
	}
	codex["zero"] = []Plural{{Key: "zero", Hint: "0"}}

	return codex
}

func equalPluralCodex(got, want map[string][]Plural) bool {
	if len(got) != len(want) {
		return false
	}

	for sourceKey, targetPlurals := range want {
		if !slices.Equal(got[sourceKey], targetPlurals) {
			return false
		}
	}

	return true
}
