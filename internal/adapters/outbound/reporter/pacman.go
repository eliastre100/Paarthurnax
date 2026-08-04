package reporter

import (
	"fmt"
	"strings"
	"sync"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const defaultTerminalWidth = 80

type Pacman struct {
	mu      sync.Mutex
	program *tea.Program
}

func NewPacmanReporter() *Pacman {
	return &Pacman{}
}

func (r *Pacman) StartHandlingChanges(documentCount uint) {
	if documentCount == 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.program != nil {
		return
	}

	r.program = tea.NewProgram(newPacmanModel(documentCount))
	go func() {
		_, _ = r.program.Run()
	}()
}

func (r *Pacman) StartHandlingDocument(document string) {
	r.send(documentStartedMsg{document: document})
}

func (r *Pacman) UpdateDocumentChanges(document string, count uint) {
	r.send(documentTotalUpdatedMsg{document: document, total: count})
}

func (r *Pacman) DoneUpdatingSegment(document string, name string) {
	r.send(segmentCompletedMsg{document: document})
}

func (r *Pacman) DoneUpdatingDocument(document string) {
	r.send(documentCompletedMsg{document: document})
}

func (r *Pacman) DoneDeletingDocument(document string) {
	r.send(documentCompletedMsg{document: document})
}

func (r *Pacman) DoneDeletingSegment(document string, name string) {
	r.send(segmentCompletedMsg{document: document})
}

func (r *Pacman) Close() {
	r.mu.Lock()
	program := r.program
	r.mu.Unlock()
	if program != nil {
		program.Quit()
		program.Wait()
	}
}

func (r *Pacman) send(msg tea.Msg) {
	r.mu.Lock()
	program := r.program
	r.mu.Unlock()
	if program != nil {
		program.Send(msg)
	}
}

type documentStartedMsg struct{ document string }

type documentTotalUpdatedMsg struct {
	document string
	total    uint
}

type segmentCompletedMsg struct{ document string }

type documentCompletedMsg struct{ document string }

type pacmanDocument struct {
	done      uint
	total     uint
	completed bool
	progress  progress.Model
}

type pacmanModel struct {
	totalDocuments uint
	doneDocuments  uint
	documents      map[string]pacmanDocument
	documentOrder  []string
	overall        progress.Model
	spinner        spinner.Model
	width          int
}

func newPacmanModel(totalDocuments uint) pacmanModel {
	return pacmanModel{
		totalDocuments: totalDocuments,
		documents:      make(map[string]pacmanDocument),
		overall:        newProgress(),
		spinner:        spinner.New(spinner.WithSpinner(spinner.MiniDot)),
		width:          defaultTerminalWidth,
	}
}

func newProgress() progress.Model {
	return progress.New(
		progress.WithoutPercentage(),
		progress.WithFillCharacters(progress.DefaultFullCharFullBlock, progress.DefaultEmptyCharBlock),
		progress.WithDefaultBlend(),
	)
}

func (m pacmanModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m pacmanModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case documentStartedMsg:
		if _, ok := m.documents[msg.document]; !ok {
			m.documents[msg.document] = pacmanDocument{progress: newProgress()}
			m.documentOrder = append(m.documentOrder, msg.document)
		}
	case documentTotalUpdatedMsg:
		document := m.documents[msg.document]
		document.total = msg.total
		commands = append(commands, document.progress.SetPercent(progressPercent(document.done, document.total)))
		m.documents[msg.document] = document
	case segmentCompletedMsg:
		document := m.documents[msg.document]
		document.done++
		commands = append(commands, document.progress.SetPercent(progressPercent(document.done, document.total)))
		m.documents[msg.document] = document
	case documentCompletedMsg:
		document := m.documents[msg.document]
		if document.completed {
			return m, nil
		}
		document.completed = true
		m.documents[msg.document] = document
		m.doneDocuments++
		commands = append(commands, m.overall.SetPercent(progressPercent(m.doneDocuments, m.totalDocuments)))
		if m.doneDocuments >= m.totalDocuments {
			return m, tea.Quit
		}
	case progress.FrameMsg:
		var command tea.Cmd
		m.overall, command = m.overall.Update(msg)
		commands = append(commands, command)
		for _, name := range m.documentOrder {
			document := m.documents[name]
			document.progress, command = document.progress.Update(msg)
			m.documents[name] = document
			commands = append(commands, command)
		}
	}

	var spinnerCommand tea.Cmd
	m.spinner, spinnerCommand = m.spinner.Update(msg)
	commands = append(commands, spinnerCommand)
	return m, tea.Batch(commands...)
}

func (m pacmanModel) View() tea.View {
	var view strings.Builder
	nameWidth := m.longestDocumentName()
	counterWidth := m.progressCounterWidth()

	for _, name := range m.documentOrder {
		document := m.documents[name]
		if document.completed {
			fmt.Fprintf(&view, "✓  %s\n", name)
			continue
		}
		namePadding := strings.Repeat(" ", max(0, nameWidth-lipgloss.Width(name)))
		prefix := fmt.Sprintf("%s  %s%s  ", m.spinner.View(), name, namePadding)
		counter := fmt.Sprintf("%d/%d", document.done, document.total)
		document.progress.SetWidth(m.progressBarWidth(lipgloss.Width(prefix), counterWidth))
		fmt.Fprintf(&view, "%s%s  %-*s\n", prefix, document.progress.View(), counterWidth, counter)
	}
	if m.doneDocuments < m.totalDocuments {
		prefix := "Total  "
		counter := fmt.Sprintf("%d/%d", m.doneDocuments, m.totalDocuments)
		m.overall.SetWidth(m.progressBarWidth(lipgloss.Width(prefix), counterWidth))
		fmt.Fprintf(&view, "\n%s%s  %-*s\n", prefix, m.overall.View(), counterWidth, counter)
	}
	return tea.NewView(view.String())
}

func (m pacmanModel) progressCounterWidth() int {
	width := lipgloss.Width(fmt.Sprintf("%d/%d", m.doneDocuments, m.totalDocuments))
	for _, document := range m.documents {
		width = max(width, lipgloss.Width(fmt.Sprintf("%d/%d", document.done, document.total)))
	}
	return width
}

func (m pacmanModel) progressBarWidth(prefixWidth, counterWidth int) int {
	return max(1, m.width-prefixWidth-2-counterWidth)
}

func (m pacmanModel) longestDocumentName() int {
	width := 0
	for _, name := range m.documentOrder {
		width = max(width, lipgloss.Width(name))
	}
	return width
}

func progressPercent(done, total uint) float64 {
	if total == 0 {
		return 0
	}
	return float64(done) / float64(total)
}
