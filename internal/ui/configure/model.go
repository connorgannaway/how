package configure

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/connorgannaway/how/internal/config"
	"github.com/connorgannaway/how/internal/ui/styles"
)

type state int

const (
	stateSelectProvider state = iota
	stateSelectModel
	stateInputAPIKey
	stateInputBaseURL
	stateCustomModel
	stateSaving
	stateDone
)

// Bubbletea model for configuration ui
type Model struct {
	config           *config.Config
	state            state
	providerList     list.Model
	modelList        list.Model
	apiKeyInput      textinput.Model
	baseURLInput     textinput.Model
	customModelInput textinput.Model
	selectedProvider string
	selectedModel    string
	err              error
	validationError  string
	validationWarning string
	hasExistingKey   bool
	width            int
	height           int
}

// Item for list
type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Ensure valid URL.
// Returns (error, warning)
func validateBaseURL(rawURL string) (error, string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err), ""
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL must use http or https"), ""
	}

	// Warn if using http (not https) for remote endpoints
	if u.Scheme == "http" && !strings.HasPrefix(u.Host, "localhost") && !strings.HasPrefix(u.Host, "127.0.0.1") {
		return nil, "⚠ Using HTTP (not HTTPS) for remote endpoint. API keys will be sent unencrypted."
	}

	return nil, ""
}

// Update model's modelList based on selected provider
func (m *Model) setupModelList() {
	models := config.ProviderModels[m.selectedProvider]
	modelItems := make([]list.Item, len(models)+1)
	for i, model := range models {
		modelItems[i] = item{title: model, desc: ""}
	}
	// Add custom model option
	modelItems[len(models)] = item{title: "Model not listed?", desc: ""}

	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false

	m.modelList = list.New(modelItems, delegate, m.width, m.height-4)
	m.modelList.Title = fmt.Sprintf("%s - Select Model", m.selectedProvider)
	m.modelList.SetShowHelp(false)
	m.modelList.SetShowStatusBar(false)
	m.modelList.SetFilteringEnabled(false)
}

func NewModel(cfg *config.Config) Model {
	// Create provider list
	providers := config.GetProviders()
	providerItems := make([]list.Item, len(providers))
	for i, p := range providers {
		providerItems[i] = item{title: p, desc: ""}
	}

	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false

	providerList := list.New(providerItems, delegate, 0, 0)
	providerList.Title = "Select AI Provider"
	providerList.SetShowHelp(false)
	providerList.SetShowStatusBar(false)
	providerList.SetFilteringEnabled(false)

	// Create API key input
	apiKeyInput := textinput.New()
	apiKeyInput.EchoMode = textinput.EchoPassword
	apiKeyInput.EchoCharacter = '•'

	// Create base URL input
	baseURLInput := textinput.New()

	// Create custom model input
	customModelInput := textinput.New()

	return Model{
		config:           cfg,
		state:            stateSelectProvider,
		providerList:     providerList,
		apiKeyInput:      apiKeyInput,
		baseURLInput:     baseURLInput,
		customModelInput: customModelInput,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.providerList.SetSize(msg.Width, msg.Height-4)
		if m.modelList.Items() != nil {
			m.modelList.SetSize(msg.Width, msg.Height-4)
		}
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateSelectProvider:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				m.state = stateDone
				return m, tea.Quit
			case "enter":
				if selected, ok := m.providerList.SelectedItem().(item); ok {
					m.selectedProvider = selected.title

					// For OpenAI-Compatible skip to custom model input
					if m.selectedProvider == config.ProviderOpenAICompatible {
						m.state = stateCustomModel
						m.customModelInput.Focus()
						return m, textinput.Blink
					}

					m.state = stateSelectModel
					m.setupModelList()
				}
				return m, nil
			}

			// Pass command to the list's update method
			var cmd tea.Cmd
			m.providerList, cmd = m.providerList.Update(msg)
			return m, cmd

		case stateSelectModel:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				m.state = stateSelectProvider
				return m, nil
			case "enter":
				if selected, ok := m.modelList.SelectedItem().(item); ok {
					if selected.title == "Model not listed?" {
						// Go to custom model input
						m.state = stateCustomModel
						m.customModelInput.Focus()
						return m, textinput.Blink
					}
					m.selectedModel = selected.title

					m.state = stateInputAPIKey
					m.apiKeyInput.Focus()

					// Check if existing key exists
					existingKey, _ := config.GetAPIKeyFromKeyring(m.selectedProvider)
					m.hasExistingKey = existingKey != ""

					return m, textinput.Blink
				}
				return m, nil
			}

			// Pass command to the list's update method
			var cmd tea.Cmd
			m.modelList, cmd = m.modelList.Update(msg)
			return m, cmd

		case stateCustomModel:
			switch msg.String() {
			case "ctrl+c", "esc":
				m.state = stateSelectModel
				m.customModelInput.SetValue("")
				return m, nil
			case "enter":
				if m.customModelInput.Value() != "" {
					m.selectedModel = m.customModelInput.Value()

					// For OpenAI-Compatible, go to base URL input
					if m.selectedProvider == config.ProviderOpenAICompatible {
						m.state = stateInputBaseURL
						m.baseURLInput.Focus()

						// Pre-fill base URL if it exists
						if m.config.BaseURL != "" {
							m.baseURLInput.SetValue(m.config.BaseURL)
						}
						return m, textinput.Blink
					}

					// For other providers, go directly to API key input
					m.state = stateInputAPIKey
					m.apiKeyInput.Focus()

					// Check if existing key exists
					existingKey, _ := config.GetAPIKeyFromKeyring(m.selectedProvider)
					m.hasExistingKey = existingKey != ""

					return m, textinput.Blink
				}
				return m, nil
			}

			// Pass command to the input's update method
			var cmd tea.Cmd
			m.customModelInput, cmd = m.customModelInput.Update(msg)
			return m, cmd

		// stateInputBaseURL reached when selected provider is "OpenAI-Compatible"
		case stateInputBaseURL:
			switch msg.String() {
			case "ctrl+c", "esc":
				m.state = stateSelectModel
				m.baseURLInput.SetValue("")
				m.validationError = ""
				m.validationWarning = ""
				return m, nil
			case "enter":
				if m.baseURLInput.Value() != "" {
					// Validate base URL
					err, warning := validateBaseURL(m.baseURLInput.Value())
					if err != nil {
						m.validationError = err.Error()
						return m, nil
					}

					m.config.BaseURL = m.baseURLInput.Value()
					m.validationError = ""
					m.validationWarning = warning
					m.state = stateInputAPIKey
					m.apiKeyInput.Focus()

					// Check if existing key exists
					existingKey, _ := config.GetAPIKeyFromKeyring(m.selectedProvider)
					m.hasExistingKey = existingKey != ""

					return m, textinput.Blink
				}
				return m, nil
			default:
				// Clear validation error when user types
				m.validationError = ""
			}

			// Pass command to the input's update method
			var cmd tea.Cmd
			m.baseURLInput, cmd = m.baseURLInput.Update(msg)
			return m, cmd

		case stateInputAPIKey:
			switch msg.String() {
			case "ctrl+c", "esc":
				// For OpenAI-Compatible, go back to base URL input
				if m.selectedProvider == config.ProviderOpenAICompatible {
					m.state = stateInputBaseURL
					m.validationWarning = ""
					m.hasExistingKey = false
					return m, nil
				}
				m.state = stateSelectModel
				m.apiKeyInput.SetValue("")
				m.validationWarning = ""
				m.hasExistingKey = false
				return m, nil
			case "enter":
				// Check if existing key exists in keyring
				existingKey, _ := config.GetAPIKeyFromKeyring(m.selectedProvider)
				hasExistingKey := existingKey != ""

				// For OpenAI-Compatible, API key is optional
				// For other providers, allow proceeding if either:
				// - User entered a new key, OR
				// - An existing key already exists (user keeping it)
				if m.selectedProvider == config.ProviderOpenAICompatible || m.apiKeyInput.Value() != "" || hasExistingKey {
					m.config.SetProvider(m.selectedProvider, m.selectedModel)
					if m.apiKeyInput.Value() != "" {
						if err := config.SetAPIKeyInKeyring(m.selectedProvider, m.apiKeyInput.Value()); err != nil {
							m.err = err
							m.state = stateDone
							return m, tea.Quit
						}
					}
					// Save config
					if err := config.Save(m.config); err != nil {
						m.err = err
						m.state = stateDone
						return m, tea.Quit
					}
					m.validationWarning = ""
					m.state = stateDone
					return m, tea.Quit
				}
				return m, nil
			}

			// Pass command to the input's update method
			var cmd tea.Cmd
			m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) View() string {
	// Show different views based on state
	switch m.state {
	case stateSelectProvider:
		return lipgloss.NewStyle().MarginTop(1).Render(m.providerList.View())

	case stateSelectModel:
		return lipgloss.NewStyle().MarginTop(1).Render(m.modelList.View())

	case stateCustomModel:
		return lipgloss.JoinVertical(
			lipgloss.Left,
			styles.InputLabelStyle.Render(fmt.Sprintf("%s - Enter Model Name:", m.selectedProvider)),
			"",
			styles.InputStyle.Render(m.customModelInput.View()),
			"",
			styles.HelpStyle.Render("enter: save • esc: back"),
		)

	case stateInputBaseURL:
		sections := []string{
			styles.InputLabelStyle.Render("Enter Base URL:"),
			"",
			styles.InputStyle.Render(m.baseURLInput.View()),
		}

		// Show validation error if present
		if m.validationError != "" {
			sections = append(sections, "")
			sections = append(sections, styles.ErrorStyle.Render("✗ "+m.validationError))
		}

		sections = append(sections, "")
		sections = append(sections, styles.HelpStyle.Render("enter: continue • esc: back"))

		return lipgloss.JoinVertical(lipgloss.Left, sections...)

	case stateInputAPIKey:
		helpText := "enter: save • esc: back"
		if m.selectedProvider == config.ProviderOpenAICompatible {
			helpText = "enter: save (leave empty if no auth required) • esc: back"
		}

		sections := []string{
			styles.InputLabelStyle.Render(fmt.Sprintf("Enter %s API Key:", m.selectedProvider)),
		}

		sections = append(sections, "")
		sections = append(sections, styles.InputStyle.Render(m.apiKeyInput.View()))
		sections = append(sections, "")
		if m.hasExistingKey {
			sections = append(sections, styles.HelpStyle.Render("Existing key found, leave blank to use"))
		}
		sections = append(sections, styles.HelpStyle.Render(helpText))

		// Show warning if present (from base URL validation)
		if m.validationWarning != "" {
			sections = append(sections, "")
			sections = append(sections, styles.WarningStyle.Render(m.validationWarning))
		}

		return lipgloss.JoinVertical(lipgloss.Left, sections...)

	case stateDone:
		if m.err != nil {
			return styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		}
		return styles.SuccessStyle.Render("✓ Configuration saved!")

	default:
		return ""
	}
}

// Check if model is in done state
func (m Model) ShouldQuit() bool {
	return m.state == stateDone
}

// Get any error that occurred
func (m Model) GetError() error {
	return m.err
}

// Start configuration UI and handle exit
func Run(cfg *config.Config) error {
	m := NewModel(cfg)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// Return any error from the model
	if finalModel, ok := finalModel.(Model); ok {
		return finalModel.GetError()
	}

	return nil
}
