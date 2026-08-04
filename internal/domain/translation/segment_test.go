package translation

import (
	"reflect"
	"testing"
)

func TestNewSegment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      string
		value    string
		expected *Segment
	}{
		{
			name:  "plain text",
			key:   "en.greeting",
			value: "Hello world",
			expected: &Segment{
				Key:    "en.greeting",
				Value:  "Hello world",
				Plural: false,
				Parts: []SegmentPart{
					TextPart{Text: "Hello world"},
				},
			},
		},
		{
			name:  "mixed text and variable placeholders",
			key:   "en.inbox.count",
			value: "Hello %{name}, you have %{count} messages.",
			expected: &Segment{
				Key:    "en.inbox.count",
				Value:  "Hello %{name}, you have %{count} messages.",
				Plural: false,
				Parts: []SegmentPart{
					TextPart{Text: "Hello "},
					VariablePart{Name: "name"},
					TextPart{Text: ", you have "},
					VariablePart{Name: "count"},
					TextPart{Text: " messages."},
				},
			},
		},
		{
			name:  "variables with messier names",
			key:   "en.inbox.count",
			value: "Hello %{name }, you have %{\t\tcount      } messages.",
			expected: &Segment{
				Key:    "en.inbox.count",
				Value:  "Hello %{name }, you have %{\t\tcount      } messages.",
				Plural: false,
				Parts: []SegmentPart{
					TextPart{Text: "Hello "},
					VariablePart{Name: "name"},
					TextPart{Text: ", you have "},
					VariablePart{Name: "count"},
					TextPart{Text: " messages."},
				},
			},
		},
		{
			name:  "plural detection",
			key:   "en.inbox.count.other",
			value: "You have %{ count } messages.",
			expected: &Segment{
				Key:    "en.inbox.count.other",
				Value:  "You have %{ count } messages.",
				Plural: true,
				Parts: []SegmentPart{
					TextPart{Text: "You have "},
					VariablePart{Name: "count"},
					TextPart{Text: " messages."},
				},
			},
		},
		{
			name:  "plural detection for other categories",
			key:   "pl.cart.items.few",
			value: "Masz %{count} wiadomosci.",
			expected: &Segment{
				Key:    "pl.cart.items.few",
				Value:  "Masz %{count} wiadomosci.",
				Plural: true,
				Parts: []SegmentPart{
					TextPart{Text: "Masz "},
					VariablePart{Name: "count"},
					TextPart{Text: " wiadomosci."},
				},
			},
		},
		{
			name:  "plural without count variable",
			key:   "en.inbox.count.other",
			value: "You have messages.",
			expected: &Segment{
				Key:    "en.inbox.count.other",
				Value:  "You have messages.",
				Plural: true,
				Parts: []SegmentPart{
					TextPart{Text: "You have messages."},
				},
			},
		},
		{
			name:  "non plural when category is not terminal key part", // This still technically be a plural, but for translation purposes, we don't want to treat it as such as it is a fixed string
			key:   "en.one.more",
			value: "One more thing",
			expected: &Segment{
				Key:    "en.one.more",
				Value:  "One more thing",
				Plural: false,
				Parts: []SegmentPart{
					TextPart{Text: "One more thing"},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewSegment(tt.key, tt.value)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("NewSegment() = %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func TestSegmentHasVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    string
		variable string
		expected bool
	}{
		{
			name:     "finds existing variable",
			value:    "Hello %{name}, you have %{count} messages.",
			variable: "name",
			expected: true,
		},
		{
			name:     "returns false when variable is missing",
			value:    "Hello %{name}",
			variable: "count",
			expected: false,
		},
		{
			name:     "matches normalized variable names",
			value:    "You have %{ count } messages.",
			variable: "count",
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			segment := NewSegment("en.test", tt.value)

			if got := segment.HasVariable(tt.variable); got != tt.expected {
				t.Fatalf("HasVariable(%q) = %t, want %t", tt.variable, got, tt.expected)
			}
		})
	}
}

func TestSegmentLeafKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "unqualified key",
			key:      "greeting",
			expected: "greeting",
		},
		{
			name:     "nested key",
			key:      "en.inbox.message",
			expected: "message",
		},
		{
			name:     "empty key",
			key:      "",
			expected: "",
		},
		{
			name:     "trailing separator",
			key:      "en.inbox.",
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			segment := Segment{Key: tt.key}
			if got := segment.LeafKey(); got != tt.expected {
				t.Fatalf("LeafKey() = %q, want %q", got, tt.expected)
			}
		})
	}
}
