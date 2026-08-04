package selection

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func prompt(question string, choices []string, input io.Reader, output io.Writer) (string, error) {
	scanner := bufio.NewScanner(input)
	validChoices := choiceSet(choices)

	for {
		if _, err := fmt.Fprintf(output, "%s: ", question); err != nil {
			return "", fmt.Errorf("writing selection prompt: %w", err)
		}
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", fmt.Errorf("reading selection: %w", err)
			}
			return "", ErrAborted
		}

		choice := strings.TrimSpace(scanner.Text())
		if _, ok := validChoices[choice]; ok {
			return choice, nil
		}
		if _, err := fmt.Fprintf(output, "Invalid selection. Choose one of [%s].\n", strings.Join(choices, ", ")); err != nil {
			return "", fmt.Errorf("writing validation error: %w", err)
		}
	}
}

func choiceSet(choices []string) map[string]struct{} {
	set := make(map[string]struct{}, len(choices))
	for _, choice := range choices {
		set[choice] = struct{}{}
	}
	return set
}
