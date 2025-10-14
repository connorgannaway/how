package styles

import "github.com/charmbracelet/lipgloss"

const (
	// Default max width for text wrapping
	DefaultMaxWidth = 80
)

var (
	Primary   = lipgloss.Color("#7C3AED")
	Secondary = lipgloss.Color("#A78BFA")
	Blue      = lipgloss.Color("#5F5FD0")
	Success   = lipgloss.Color("#10B981")
	Error     = lipgloss.Color("#EF4444")
	Warning   = lipgloss.Color("#F59E0B")
	Muted     = lipgloss.Color("#9CA3AF")
	White     = lipgloss.Color("#FFFFFF")
	Black     = lipgloss.Color("#000000")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	DescriptionStyle = lipgloss.NewStyle().
				Foreground(Muted).
				MarginBottom(1)

	CommandStyle = lipgloss.NewStyle().
			Foreground(White).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)

	PromptSymbol = lipgloss.NewStyle().
				Foreground(Success).
				Render("$ ")

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	QuestionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	// -- List styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true).
				PaddingLeft(2)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(White).
			PaddingLeft(4)
	// -- 

	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted).
			MarginTop(1).
			MarginLeft(2)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning).
			MarginTop(1).
			MarginLeft(2)

	// -- Input styles
	InputStyle = lipgloss.NewStyle().
			Foreground(White).
			MarginLeft(2)

	InputLabelStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(Blue).
			Bold(false).
			MarginTop(1).
			MarginBottom(1).
			MarginLeft(2).
		 	Padding(0, 1)
)
