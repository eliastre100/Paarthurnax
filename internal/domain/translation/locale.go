package translation

import (
	"strconv"
)

//go:generate go run ./cmd/plurals -output plurals_generated.go

type Locale string

const (
	Afrikaans        Locale = "af"
	Arabic           Locale = "ar"
	Azerbaijani      Locale = "az"
	Belarusian       Locale = "be"
	Bulgarian        Locale = "bg"
	Bengali          Locale = "bn"
	Bosnian          Locale = "bs"
	Catalan          Locale = "ca"
	Montenegrin      Locale = "cnr"
	Czech            Locale = "cs"
	Welsh            Locale = "cy"
	Danish           Locale = "da"
	German           Locale = "de"
	Dzongkha         Locale = "dz"
	Greek            Locale = "el"
	English          Locale = "en"
	Esperanto        Locale = "eo"
	Spanish          Locale = "es"
	Estonian         Locale = "et"
	Basque           Locale = "eu"
	Persian          Locale = "fa"
	Finnish          Locale = "fi"
	French           Locale = "fr"
	WesternFrisian   Locale = "fy"
	Gaelic           Locale = "gd"
	Galician         Locale = "gl"
	Hebrew           Locale = "he"
	Hindi            Locale = "hi"
	Croatian         Locale = "hr"
	Hungarian        Locale = "hu"
	Armenian         Locale = "hy"
	Indonesian       Locale = "id"
	Icelandic        Locale = "is"
	Italian          Locale = "it"
	Japanese         Locale = "ja"
	Georgian         Locale = "ka"
	Kazakh           Locale = "kk"
	CentralKhmer     Locale = "km"
	Kannada          Locale = "kn"
	Korean           Locale = "ko"
	Luxembourgish    Locale = "lb"
	Lao              Locale = "lo"
	Lithuanian       Locale = "lt"
	Latvian          Locale = "lv"
	Malagasy         Locale = "mg"
	Macedonian       Locale = "mk"
	Malayalam        Locale = "ml"
	Mongolian        Locale = "mn"
	Malay            Locale = "ms"
	NorwegianBokmål  Locale = "nb"
	Nepali           Locale = "ne"
	Dutch            Locale = "nl"
	NorwegianNynorsk Locale = "nn"
	Occitan          Locale = "oc"
	Odia             Locale = "or"
	Punjabi          Locale = "pa"
	Polish           Locale = "pl"
	Portuguese       Locale = "pt"
	Romansh          Locale = "rm"
	Romanian         Locale = "ro"
	Russian          Locale = "ru"
	Sardinian        Locale = "sc"
	Slovak           Locale = "sk"
	Slovenian        Locale = "sl"
	Albanian         Locale = "sq"
	Serbian          Locale = "sr"
	SouthernSotho    Locale = "st"
	Swedish          Locale = "sv"
	Swahili          Locale = "sw"
	Tamil            Locale = "ta"
	Telugu           Locale = "te"
	Thai             Locale = "th"
	Tagalog          Locale = "tl"
	Turkish          Locale = "tr"
	Tatar            Locale = "tt"
	Uyghur           Locale = "ug"
	Ukrainian        Locale = "uk"
	Urdu             Locale = "ur"
	Uzbek            Locale = "uz"
	Vietnamese       Locale = "vi"
	Wolof            Locale = "wo"
)

var Locales = []Locale{
	Afrikaans,
	Arabic,
	Azerbaijani,
	Belarusian,
	Bulgarian,
	Bengali,
	Bosnian,
	Catalan,
	Montenegrin,
	Czech,
	Welsh,
	Danish,
	German,
	Dzongkha,
	Greek,
	English,
	Esperanto,
	Spanish,
	Estonian,
	Basque,
	Persian,
	Finnish,
	French,
	WesternFrisian,
	Gaelic,
	Galician,
	Hebrew,
	Hindi,
	Croatian,
	Hungarian,
	Armenian,
	Indonesian,
	Icelandic,
	Italian,
	Japanese,
	Georgian,
	Kazakh,
	CentralKhmer,
	Kannada,
	Korean,
	Luxembourgish,
	Lao,
	Lithuanian,
	Latvian,
	Malagasy,
	Macedonian,
	Malayalam,
	Mongolian,
	Malay,
	NorwegianBokmål,
	Nepali,
	Dutch,
	NorwegianNynorsk,
	Occitan,
	Odia,
	Punjabi,
	Polish,
	Portuguese,
	Romansh,
	Romanian,
	Russian,
	Sardinian,
	Slovak,
	Slovenian,
	Albanian,
	Serbian,
	SouthernSotho,
	Swedish,
	Swahili,
	Tamil,
	Telugu,
	Thai,
	Tagalog,
	Turkish,
	Tatar,
	Uyghur,
	Ukrainian,
	Urdu,
	Uzbek,
	Vietnamese,
	Wolof,
}

type Plural struct {
	Key  string
	Hint string
}

func (l Locale) String() string {
	return string(l)
}

// PluralCodex returns the destination plural keys affected by each source plural key.
func (l Locale) PluralCodex(to Locale) map[string][]Plural {
	codex := make(map[string][]Plural)
	source, ok := plurals[l.String()]
	if !ok {
		return codex
	}
	target, ok := plurals[to.String()]
	if !ok {
		return codex
	}

	for _, hint := range target.Hints {
		if hint.Key == "zero" {
			continue
		}

		key := source.Evaluator(hint.Hint)
		if _, ok := codex[string(key)]; !ok {
			codex[string(key)] = []Plural{}
		}
		codex[string(key)] = append(codex[string(key)], Plural{Key: hint.Key, Hint: strconv.Itoa(hint.Hint)})
	}

	if _, ok := codex["zero"]; !ok {
		codex["zero"] = []Plural{}
	}
	codex["zero"] = append(codex["zero"], Plural{Key: "zero", Hint: "0"})

	return codex
}
