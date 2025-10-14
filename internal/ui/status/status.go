package status

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/connorgannaway/how/internal/config"
	"github.com/connorgannaway/how/internal/ui/styles"
)

var (
	labelStyle = lipgloss.NewStyle().
			Foreground(styles.Primary).
			Bold(true).
			Width(10).
			Align(lipgloss.Right)

	valueStyle = lipgloss.NewStyle().
			Foreground(styles.White)

	providerItemStyle = lipgloss.NewStyle().
				Foreground(styles.Secondary).
				Width(19).
				Align(lipgloss.Right).
				Bold(true)

	notSetStyle = lipgloss.NewStyle().
			Foreground(styles.Muted).
			Italic(true)
)

// Format api key for printing
func retrieveAndFormatKey(provider string, revealFull bool) string {
	var keyLine string
	apiKey, _ := config.GetAPIKeyFromKeyring(provider)
	if apiKey != "" {
		if !revealFull {
			apiKey = config.MaskAPIKey(apiKey)
		}
		keyLine = valueStyle.Render(apiKey)
	} else {
		keyLine = notSetStyle.Render("(not set)")
	}

	return keyLine
}

// Run displays the configuration status with styled output.
// This is not a bubbletea model
func Run(cfg *config.Config, showKey, showAll, revealFull bool) {
	// Test for valid configuration
	if ready, missing := cfg.IsConfigured(); !ready {
		output := lipgloss.JoinVertical(
			lipgloss.Left,
			styles.ErrorStyle.Render("Status: Not configured"),
			styles.MutedStyle.Render(fmt.Sprintf("Missing: %v", strings.Join(missing, ", "))),
			styles.MutedStyle.Render("Run 'how --configure' to set up."),
		)
		fmt.Println("\n" + lipgloss.NewStyle().Padding(0, 2).Render(output))
		return
	}

	var lines []string

	// Provider and Model section
	providerLine := fmt.Sprintf("%s %s",
		labelStyle.Render("Provider:"),
		valueStyle.Render(cfg.CurrentProvider),
	)
	modelLine := fmt.Sprintf("%s %s",
		labelStyle.Render("Model:"),
		valueStyle.Render(cfg.CurrentModel),
	)

	lines = append(lines, providerLine, modelLine)

	// Base URL for OpenAI-Compatible provider
	if cfg.CurrentProvider == config.ProviderOpenAICompatible {
		baseURLLine := fmt.Sprintf("%s %s",
			labelStyle.Render("Base URL:"),
			valueStyle.Render(cfg.BaseURL),
		)
		lines = append(lines, baseURLLine)
	}

	// Show API key(s)
	if showKey {
		if showAll {
			lines = append(lines, labelStyle.Render("\nAPI Keys:"))

			//Create a line for each provider
			providers := config.GetProviders()
			for _, provider := range providers {
				keyValue := retrieveAndFormatKey(provider, revealFull)
				apiKeyLine := fmt.Sprintf("%s %s",
					providerItemStyle.Render(provider+":"),
					keyValue,
				)
				
				lines = append(lines, apiKeyLine)
				}
			
		} else {
			// Show only current provider's API key
			keyValue := retrieveAndFormatKey(cfg.CurrentProvider, revealFull)
			apiKeyLine := fmt.Sprintf("%s %s",
				labelStyle.Render("API Key:"),
				keyValue,
			)
			
			lines = append(lines, apiKeyLine)
		}
	}

	output := strings.Join(lines, "\n")
	fmt.Println("\n" + lipgloss.NewStyle().Padding(0, 2).Render(output))
}
