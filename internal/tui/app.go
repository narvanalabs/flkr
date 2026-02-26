package tui

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/narvanalabs/flkr/internal/generator"
	"github.com/narvanalabs/flkr/pkg/flkr"
)

// step tracks the current phase of the TUI wizard.
type step int

const (
	stepDetect  step = iota
	stepReview
	stepConfirm
	stepDone
)

// Model is the top-level Bubble Tea model for the flkr init wizard.
type Model struct {
	path            string
	templateVersion string

	step    step
	spinner spinner.Model
	profile *flkr.AppProfile
	err     error

	reviewForm  *huh.Form
	confirmForm *huh.Form
	portStr     *string
	confirmed   *bool
	preview     string
	outputPath  string
}

// New creates a new TUI model.
func New(path, templateVersion string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	confirmed := false

	return Model{
		path:            path,
		templateVersion: templateVersion,
		step:            stepDetect,
		spinner:         s,
		confirmed:       &confirmed,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, runDetection(m.path))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	switch m.step {
	case stepDetect:
		return m.updateDetect(msg)
	case stepReview:
		return m.updateReview(msg)
	case stepConfirm:
		return m.updateConfirm(msg)
	case stepDone:
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
	header := titleStyle.Render("flkr init") + "\n"

	switch m.step {
	case stepDetect:
		if m.profile != nil {
			return header + renderDetectResult(m.profile)
		}
		if m.err != nil {
			return header + errorStyle.Render(fmt.Sprintf("  Error: %v", m.err)) + "\n"
		}
		return header + renderDetecting(m.spinner)

	case stepReview:
		s := header + subtitleStyle.Render("  Review detected settings:") + "\n\n"
		if m.reviewForm != nil {
			s += m.reviewForm.View()
		}
		return s

	case stepConfirm:
		s := header + subtitleStyle.Render("  Preview:") + "\n\n"
		s += borderStyle.Render(m.preview) + "\n\n"
		if m.confirmForm != nil {
			s += m.confirmForm.View()
		}
		return s

	case stepDone:
		if m.confirmed != nil && *m.confirmed {
			return header + successStyle.Render(fmt.Sprintf("  Wrote %s", m.outputPath)) + "\n\n"
		}
		return header + dimStyle.Render("  Cancelled.") + "\n\n"
	}

	return header
}

func (m Model) updateDetect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case detectResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.step = stepDone
			return m, tea.Quit
		}
		if msg.profile == nil {
			m.err = fmt.Errorf("no application stack detected")
			m.step = stepDone
			return m, tea.Quit
		}
		m.profile = msg.profile
		portStr := strconv.Itoa(m.profile.Port)
		m.portStr = &portStr
		m.reviewForm = buildReviewForm(m.profile, m.portStr)
		m.step = stepReview
		return m, m.reviewForm.Init()

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Model) updateReview(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.reviewForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.reviewForm = f
	}

	if m.reviewForm.State == huh.StateCompleted {
		applyFormValues(m.profile, *m.portStr)

		preview, err := generatePreview(m.profile, m.templateVersion)
		if err != nil {
			m.err = err
			m.step = stepDone
			return m, tea.Quit
		}
		m.preview = preview
		m.confirmForm = buildConfirmForm(m.confirmed)
		m.step = stepConfirm
		return m, m.confirmForm.Init()
	}

	if m.reviewForm.State == huh.StateAborted {
		m.step = stepDone
		return m, tea.Quit
	}

	return m, cmd
}

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.confirmForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.confirmForm = f
	}

	if m.confirmForm.State == huh.StateCompleted {
		if *m.confirmed {
			m.outputPath = filepath.Join(m.path, "flake.nix")
			gen := &generator.DefaultGenerator{}
			_, err := gen.Generate(m.profile, generator.Options{
				OutputPath:      m.outputPath,
				TemplateVersion: m.templateVersion,
			})
			if err != nil {
				m.err = err
			}
		}
		m.step = stepDone
		return m, tea.Quit
	}

	if m.confirmForm.State == huh.StateAborted {
		m.step = stepDone
		return m, tea.Quit
	}

	return m, cmd
}
