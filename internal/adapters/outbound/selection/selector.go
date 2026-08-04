package selection

import (
	"errors"
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
)

var (
	ErrNoChoices = errors.New("select requires at least one choice")
	ErrAborted   = errors.New("selection aborted")
)

const defaultVisibleChoices = 10

type Selector struct {
	input       io.Reader
	output      io.Writer
	interactive bool
}

func NewSelector() *Selector {
	return newSelector(os.Stdin, os.Stdout, hasTerminal())
}

func newSelector(input io.Reader, output io.Writer, interactive bool) *Selector {
	return &Selector{
		input:       input,
		output:      output,
		interactive: interactive,
	}
}

func (s *Selector) Select(question string, choices []string) (string, error) {
	return execute(question, choices, s.input, s.output, s.interactive)
}

func execute(question string, choices []string, input io.Reader, output io.Writer, interactive bool) (string, error) {
	if len(choices) == 0 {
		return "", ErrNoChoices
	}

	if !interactive {
		return prompt(question, choices, input, output)
	}

	choice, err := runTUI(question, choices, input, output)
	if err != nil {
		return "", writeCancellationResult(output, question, err)
	}
	if err := writeResult(output, question, choice); err != nil {
		return "", err
	}
	return choice, nil
}

func writeCancellationResult(output io.Writer, question string, err error) error {
	if !errors.Is(err, ErrAborted) {
		return err
	}
	if writeErr := writeResult(output, question, "<canceled>"); writeErr != nil {
		return writeErr
	}
	return err
}

func writeResult(output io.Writer, question, result string) error {
	if _, err := fmt.Fprintf(output, "%s: %s\n", question, result); err != nil {
		return fmt.Errorf("writing selection: %w", err)
	}
	return nil
}

func hasTerminal() bool {
	return term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd())
}

func runTUI(question string, choices []string, input io.Reader, output io.Writer) (string, error) {
	program := tea.NewProgram(newSelectionModel(question, choices), tea.WithInput(input), tea.WithOutput(output))
	result, err := program.Run()
	if err != nil {
		return "", fmt.Errorf("running selection interface: %w", err)
	}

	model, ok := result.(selectionModel)
	if !ok {
		return "", errors.New("unexpected selection interface result")
	}
	return model.result()
}
