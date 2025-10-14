package question

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/connorgannaway/how/internal/ai"
	"github.com/connorgannaway/how/internal/clipboard"
	"github.com/connorgannaway/how/internal/system"
	"github.com/connorgannaway/how/internal/ui/styles"
)

type state int

const (
	stateThinking state = iota
	stateDisplaying
	stateError
	stateDone
)

// Bubbletea model for question UI
type Model struct {
	question string
	provider ai.Provider
	sysInfo  *system.SystemInfo
	spinner  spinner.Model
	state    state
	response *ai.Response
	err      error
	copied   bool
	width    int
}

// Bubbletea messages
type aiResponseMsg struct {
	response *ai.Response
}

type aiErrorMsg struct {
	err error
}

type clipboardMsg struct {
	success bool
}

// wrapper for ai.Provider.Ask to usage with model and tea commands
func (m Model) askAI() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		response, err := m.provider.Ask(ctx, m.question, m.sysInfo)
		if err != nil {
			return aiErrorMsg{err: err}
		}
		return aiResponseMsg{response: response}
	}
}

// clipboard wrapper for usage with tea commands
func copyToClipboard(commands []string) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.CopyCommands(commands)
		return clipboardMsg{success: err == nil}
	}
}

func NewModel(question string, provider ai.Provider, sysInfo *system.SystemInfo) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle

	return Model{
		question: question,
		provider: provider,
		sysInfo:  sysInfo,
		spinner:  s,
		state:    stateThinking,
	}
}

// Kick off spinner and send question to AI provider
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.askAI(),
	)
}

// Handle messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.state = stateDone
			return m, tea.Quit
		case "enter":
			if m.state == stateDisplaying || m.state == stateError {
				m.state = stateDone
				return m, tea.Quit
			}
		}

	case aiResponseMsg:
		m.response = msg.response
		m.state = stateDisplaying
		// Auto-copy to clipboard
		if len(m.response.Commands) > 0 {
			return m, tea.Sequence(copyToClipboard(m.response.Commands), tea.Quit)
		}
		return m, tea.Quit

	case aiErrorMsg:
		m.err = msg.err
		m.state = stateError
		return m, tea.Quit

	case clipboardMsg:
		m.copied = msg.success
		return m, nil

	case spinner.TickMsg:
		if m.state == stateThinking {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) View() string {
    var parts []string

    // Calculate effective width for text wrapping (min of terminal width - padding, or max 80)
    effectiveWidth := styles.DefaultMaxWidth
    if m.width > 0 {
        availableWidth := m.width - 4
        if availableWidth < effectiveWidth {
            effectiveWidth = availableWidth
        }
    }

    // Show question
    parts = append(parts, styles.QuestionStyle.Width(effectiveWidth).Render("⚡ " + ai.BuildUserPrompt(m.question)))
    parts = append(parts, "")

    switch m.state {

	// Show spinner while waiting for response
    case stateThinking:
        parts = append(parts, m.spinner.View()+" "+styles.MutedStyle.Width(effectiveWidth).Render("Thinking..."))

	// Display response parts
    case stateDisplaying:
        if m.response != nil {
            if m.response.Title != "" {
                parts = append(parts, styles.TitleStyle.Width(effectiveWidth).Render(m.response.Title))
            }
            if m.response.Description != "" {
                parts = append(parts, styles.DescriptionStyle.Width(effectiveWidth).Render(m.response.Description))
            }
            if len(m.response.Commands) > 0 {
                for _, cmd := range m.response.Commands {
					// Don't render prompt symbol for potential scripts
                    if strings.Contains(cmd, "\n") {
                        parts = append(parts, styles.CommandStyle.Width(effectiveWidth).Render(cmd))
                    } else {
                        parts = append(parts, styles.CommandStyle.Width(effectiveWidth).Render(styles.PromptSymbol + cmd))
                    }
                }
            }
            if m.copied {
                parts = append(parts, "", styles.SuccessStyle.Render("✓ Copied to clipboard"))
            }
        }

    case stateError:
        errorText := fmt.Sprintf("Error: %v", m.err)
        parts = append(parts, styles.ErrorStyle.Width(effectiveWidth).Render(errorText))
    }

    return lipgloss.NewStyle().Padding(1, 2).Render(strings.Join(parts, "\n"))
}

// Check if model is in the done state
func (m Model) ShouldQuit() bool {
	return m.state == stateDone
}

// GetError returns any error that occurred
func (m Model) GetError() error {
	return m.err
}

// Start UI and handle exit
func Run(question string, provider ai.Provider, sysInfo *system.SystemInfo) error {
	m := NewModel(question, provider, sysInfo)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// If there was an error, we already displayed it in the UI, so return nil
	// to prevent main.go from printing it again
	if finalModel, ok := finalModel.(Model); ok {
		if finalModel.GetError() != nil {
			os.Exit(1)
		}
	}

	return nil
}
