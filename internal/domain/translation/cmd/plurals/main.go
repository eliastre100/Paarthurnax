package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"text/template"
)

type pluralizationRule struct {
	Keys     []string `json:"keys"`
	RuleName string   `json:"rule_name"`
}

type pluralHints map[string]map[string]int

type generatedHint struct {
	Key   string
	Value int
}

type generatedDefinition struct {
	Evaluator string
	RuleName  string
	Hints     []generatedHint
}

//go:embed evaluators/*.go.tmpl
var evaluatorTemplates embed.FS

//go:embed hints.json
var hintsSource []byte

//go:embed plurals.go.tmpl
var pluralsTemplate string

func main() {
	output := flag.String("output", "plurals_generated.go", "generated Go file")
	flag.Parse()

	hints, err := loadHints()
	if err != nil {
		fail(err)
	}

	rules, err := extractPluralizationRules()
	if err != nil {
		fail(err)
	}

	definitions, err := buildDefinitions(rules, hints)
	if err != nil {
		fail(err)
	}

	if err := writeGenerated(*output, definitions); err != nil {
		fail(err)
	}
}

func loadHints() (pluralHints, error) {
	var hints pluralHints
	if err := json.Unmarshal(hintsSource, &hints); err != nil {
		return nil, fmt.Errorf("decode plural hints: %w", err)
	}

	return hints, nil
}

func extractPluralizationRules() (map[string]pluralizationRule, error) {
	dir, err := rubyDirectory()
	if err != nil {
		return nil, err
	}

	if output, err := run(dir, "bundle", "install"); err != nil {
		return nil, fmt.Errorf("install Ruby gems with Bundler: %w\n%s", err, output)
	}

	output, err := run(dir, "bundle", "exec", "ruby", "extract_pluralization_rules.rb")
	if err != nil {
		return nil, fmt.Errorf("extract pluralization rules: %w\n%s", err, output)
	}

	var rules map[string]pluralizationRule
	if err := json.Unmarshal(output, &rules); err != nil {
		return nil, fmt.Errorf("decode pluralization rules: %w", err)
	}

	return rules, nil
}

func rubyDirectory() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("locate plural generator")
	}

	return filepath.Join(filepath.Dir(file), "ruby"), nil
}

func run(dir, name string, args ...string) ([]byte, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, fmt.Errorf("%s is required: %w", name, err)
	}

	command := exec.Command(path, args...)
	command.Dir = dir
	return command.CombinedOutput()
}

func buildDefinitions(rules map[string]pluralizationRule, hints pluralHints) (map[string]generatedDefinition, error) {
	definitions := make(map[string]generatedDefinition, len(rules))
	for locale, rule := range rules {
		ruleName := strings.ToLower(rule.RuleName)
		ruleHints, ok := hints[ruleName]
		if !ok {
			return nil, fmt.Errorf("%s uses pluralization rule %q without hints", locale, rule.RuleName)
		}

		definition := generatedDefinition{
			Evaluator: evaluatorName(ruleName),
			RuleName:  ruleName,
			Hints:     make([]generatedHint, 0, len(rule.Keys)),
		}
		for _, key := range rule.Keys {
			hint, ok := ruleHints[key]
			if !ok {
				return nil, fmt.Errorf("%s rule %q has no hint for %q", locale, rule.RuleName, key)
			}
			definition.Hints = append(definition.Hints, generatedHint{Key: key, Value: hint})
		}
		definitions[locale] = definition
	}

	return definitions, nil
}

func writeGenerated(path string, definitions map[string]generatedDefinition) error {
	rulesSource, err := readEvaluatorTemplates(definitions)
	if err != nil {
		return err
	}

	tmpl, err := template.New("plurals").Parse(pluralsTemplate)
	if err != nil {
		return fmt.Errorf("parse plural template: %w", err)
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, struct {
		Package     string
		Definitions map[string]generatedDefinition
		Rules       string
	}{
		Package:     "translation",
		Definitions: definitions,
		Rules:       rulesSource,
	})
	if err != nil {
		return fmt.Errorf("render plural template: %w", err)
	}

	formatted, err := format.Source(output.Bytes())
	if err != nil {
		return fmt.Errorf("format generated plural source: %w", err)
	}

	return os.WriteFile(path, formatted, 0o644)
}

func evaluatorName(ruleName string) string {
	parts := strings.Split(ruleName, "_")
	for i := range parts {
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}

	return strings.Join(parts, "") + "Rule"
}

func readEvaluatorTemplates(definitions map[string]generatedDefinition) (string, error) {
	names := make([]string, 0, len(definitions))
	for _, definition := range definitions {
		names = append(names, definition.RuleName)
	}
	sort.Strings(names)
	names = slices.Compact(names)
	names = append([]string{"common"}, names...)

	var source bytes.Buffer
	for _, name := range names {
		content, err := evaluatorTemplates.ReadFile("evaluators/" + name + ".go.tmpl")
		if err != nil {
			return "", fmt.Errorf("read evaluator template for %q: %w", name, err)
		}
		if source.Len() > 0 {
			source.WriteByte('\n')
		}
		source.Write(content)
	}

	return source.String(), nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
