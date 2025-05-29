package translationgroup

var keysDirections = map[string]map[string][]pluralDefinition{
	"zero": {
		"en": {
			{
				key: "zero",
				tip: 0,
			},
		},
		"es": {
			{
				key: "zero",
				tip: 0,
			},
		},
		"it": {
			{
				key: "zero",
				tip: 0,
			},
		},
		"de": {
			{
				key: "zero",
				tip: 0,
			},
		},
		"hu": {
			{
				key: "zero",
				tip: 0,
			},
		},
		"pt": {
			{
				key: "zero",
				tip: 0,
			},
		},
		"pl": {
			{
				key: "zero",
				tip: 0,
			},
		},
		"ro": {
			{
				key: "zero",
				tip: 0,
			},
		},
		"uk": {
			{
				key: "zero",
				tip: 0,
			},
		},
	},
	"one": {
		"en": {
			{
				key: "one",
				tip: 1,
			},
		},
		"es": {
			{
				key: "one",
				tip: 1,
			},
		},
		"it": {
			{
				key: "one",
				tip: 1,
			},
		},
		"de": {
			{
				key: "one",
				tip: 1,
			},
		},
		"hu": {
			{
				key: "one",
				tip: 1,
			},
		},
		"pt": {
			{
				key: "one",
				tip: 1,
			},
		},
		"pl": {
			{
				key: "one",
				tip: 1,
			},
		},
		"ro": {
			{
				key: "one",
				tip: 1,
			},
		},
		"uk": {
			{
				key: "one",
				tip: 1,
			},
		},
	},
	"other": {
		"en": {
			{
				key: "other",
				tip: 2,
			},
		},
		"es": {
			{
				key: "other",
				tip: 2,
			},
		},
		"it": {
			{
				key: "other",
				tip: 2,
			},
		},
		"de": {
			{
				key: "other",
				tip: 2,
			},
		},
		"hu": {
			{
				key: "other",
				tip: 2,
			},
		},
		"pt": {
			{
				key: "other",
				tip: 2,
			},
		},
		"pl": {
			{
				key: "few",
				tip: 2,
			},
			{
				key: "many",
				tip: 11,
			},
			{
				key: "other",
				tip: 2, // No value match
			},
		}, // https://github.com/svenfuchs/rails-i18n/blob/master/rails/pluralization/pl.rb
		"ro": {
			{
				key: "few",
				tip: 2,
			},
			{
				key: "other",
				tip: 20,
			},
		}, // https://github.com/svenfuchs/rails-i18n/blob/master/lib/rails_i18n/common_pluralizations/romanian.rb#L5
		"uk": {
			{
				key: "few",
				tip: 2,
			},
			{
				key: "many",
				tip: 20,
			},
			{
				key: "other",
				tip: 11,
			},
		}, // https://github.com/svenfuchs/rails-i18n/blob/master/lib/rails_i18n/common_pluralizations/east_slavic.rb#L8
	},
}
