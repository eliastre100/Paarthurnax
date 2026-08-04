# Paarthurnax

Paarthurnax keeps Rails-style YAML locale files in sync by translating source
locale changes with the [DeepL API](https://www.deepl.com/docs-api).

## Requirements

- Go 1.26.4 or later, to build from source.
- A DeepL **Free** API key. Paid DeepL API keys are not currently supported.
- Translation files written in YAML. Paarthurnax handles `.yml` and `.yaml`
  files only; other translation formats are not supported.

Build the CLI from the repository root:

```sh
go build -o paarthurnax .
```

Make your DeepL Free API key available to the CLI:

```sh
export DEEPL_API_KEY='your-deepl-free-api-key'
```

## Locale Files

Paarthurnax scans only `config/locales` and its subdirectories for YAML files.
Each file must use a top-level locale key with nested translation keys and
string values, as in a standard Rails locale file:

```yaml
en:
  greeting: Hello!
  account:
    welcome: Welcome, %{name}!
```

Use separate locale files such as `config/locales/en.yml` and
`config/locales/fr.yml`, or a YAML file containing multiple top-level locales.
Values must be strings (or nested maps ending in strings); arrays, numbers,
booleans, and other YAML value types are not handled.

## Usage

Run commands from the root of the application whose locale files you want to
manage.

1. Initialise Paarthurnax after your source locale files are in their desired
   starting state:

   ```sh
   paarthurnax init
   ```

   If more than one locale is present, select the source locale when prompted.
   This creates `.paarthurnax`, which records the source-language snapshot and
   project settings. Commit this file with the project.

2. Add each destination locale:

   ```sh
   paarthurnax locales add
   ```

   Choose a locale from the interactive prompt. Paarthurnax translates all
   existing source strings and creates or fills the corresponding destination
   locale files.

3. After adding, changing, or removing source strings, run:

   ```sh
   paarthurnax translate
   ```

   Paarthurnax compares the source locale with its saved snapshot, translates
   new or changed strings into every configured destination locale, and removes
   translations for deleted source strings.

*Optional:* Rewrite locale files into a consistent normalized form to reduce
ordering noise:

   ```sh
   paarthurnax normalize
   ```

Use `paarthurnax --help` or `paarthurnax <command> --help` to inspect the
available commands.

## Current Limitations

- Only YAML locale files (`.yml` and `.yaml`) are supported.
- Only DeepL Free API keys are supported; the CLI always uses the Free API
  endpoint.
- Locale configuration is interactive.
- Review machine translations before shipping them.

## Development

Run the test suite with:

```sh
go test ./...
```

### Updating Plural Rules

Pluralization rules are generated from `rails-i18n`. To refresh them, first
update the `rails-i18n` dependency in
`internal/domain/translation/cmd/plurals/ruby/Gemfile.lock` (
`bundle update rails-i18n`).

Then regenerate the Go source from the repository root:

```sh
go generate ./internal/domain/translation
```

This requires Ruby and Bundler. 
