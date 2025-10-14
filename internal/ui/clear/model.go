package clear

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/connorgannaway/how/internal/config"
	"github.com/connorgannaway/how/internal/ui/styles"
)

type state int

const (
	stateSelecting state = iota
	stateDone
)

type providerItem struct {
	name      string
	hasKey    bool
	selected  bool
}

// Bubbletea model for clearing API keys
type Model struct {
	config    *config.Config
	state     state
	providers []providerItem
	cursor    int
	cleared   int
	err       error
}

// Create model for key clearing UI
func NewModel(cfg *config.Config) Model {
	allProviders := config.GetProviders()
	items := make([]providerItem, len(allProviders))

	for i, provider := range allProviders {
		apiKey, _ := config.GetAPIKeyFromKeyring(provider)
		items[i] = providerItem{
			name:     provider,
			hasKey:   apiKey != "",
			selected: false,
		}
	}

	return Model{
		config:    cfg,
		state:     stateSelecting,
		providers: items,
		cursor:    0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == stateSelecting {
			switch msg.String() {
			// quit
			case "ctrl+c", "q", "esc":
				m.state = stateDone
				return m, tea.Quit

			// list controls
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.providers)-1 {
					m.cursor++
				}
			case " ":
				m.providers[m.cursor].selected = !m.providers[m.cursor].selected

			// Clear selected keys
			case "enter":
				for _, provider := range m.providers {
					if provider.selected && provider.hasKey {
						if err := config.DeleteAPIKeyFromKeyring(provider.name); err != nil {
							m.err = err
							m.state = stateDone
							return m, tea.Quit
						}
						m.cleared++
					}
				}

				m.state = stateDone
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateSelecting:
		var items []string

		// Header
		title := styles.InputLabelStyle.Render("Select API keys to clear:")
		items = append(items, title)
		items = append(items, "")

		// List
		for i, provider := range m.providers {
			checkbox := "[ ]"
			if provider.selected {
				checkbox = "[x]"
			}

			cursor := "  "
			if m.cursor == i {
				cursor = styles.SuccessStyle.Render("> ")
			}

			keyStatus := ""
			if provider.hasKey {
				keyStatus = styles.SuccessStyle.Render(" (has key)")
			} else {
				keyStatus = styles.MutedStyle.Render(" (no key)")
			}

			line := fmt.Sprintf("%s%s %s%s", cursor, checkbox, provider.name, keyStatus)
			items = append(items, line)
		}

		// Footer
		items = append(items, "")
		items = append(items, styles.HelpStyle.Render("space: toggle • enter: clear selected • esc: cancel"))

		return lipgloss.NewStyle().Padding(1, 2).Render(lipgloss.JoinVertical(lipgloss.Left, items...))

	case stateDone:
		return ""

	default:
		return ""
	}
}

// ShouldQuit returns true if the model is done
func (m Model) ShouldQuit() bool {
	return m.state == stateDone
}

// GetError returns any error that occurred
func (m Model) GetError() error {
	return m.err
}

// Run starts the clear UI
func Run(cfg *config.Config) error {
	m := NewModel(cfg)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	if finalModel, ok := finalModel.(Model); ok {
		if finalModel.GetError() != nil {
			return finalModel.GetError()
		}

		// Print confirmation message
		if finalModel.cleared > 0 {
			fmt.Printf("✓ Cleared %d API key(s)\n", finalModel.cleared)
		} else {
			fmt.Println("No API keys cleared.")
		}
	}

	return nil
}
