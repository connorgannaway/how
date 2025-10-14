package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/connorgannaway/how/internal/system"
)

type Provider interface {
	// Sends a question with system context and returns the AI response
	Ask(ctx context.Context, question string, sysInfo *system.SystemInfo) (*Response, error)
	GetName() string
}

type Response struct {
	Title       string   // Optional 1-liner title
	Description string   // Optional description
	Commands    []string // commands or script lines
	RawResponse string   // Raw AI response text
}

// ParseResponse parses an AI response string into a Response struct
func ParseResponse(rawResponse string) *Response {
	response := &Response{
		RawResponse: rawResponse,
		Commands:    make([]string, 0),
	}

	lines := strings.Split(rawResponse, "\n")

	var passedScriptMarker bool
	var scriptLines []string
	var awaitingTitle bool
	var awaitingDescription bool
	var awaitingCommand bool

	// Iterate through lines to find markers
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "TITLE:") {
			response.Title = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "TITLE:"))
			if response.Title != "" {
				awaitingTitle = false
				awaitingDescription = false
				awaitingCommand = false
			} else {
				// Title is on the next line
				awaitingTitle = true
			}
			continue
		}

		if strings.HasPrefix(trimmedLine, "DESCRIPTION:") {
			response.Description = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "DESCRIPTION:"))
			if response.Description == "" {
				awaitingTitle = false
				awaitingDescription = false
				awaitingCommand = false
			} else {
				// Description is on the next line
				awaitingDescription = true
			}
			continue
		}

		if strings.HasPrefix(trimmedLine, "COMMAND:") {
			cmd := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "COMMAND:"))
			if cmd != "" {
				response.Commands = append(response.Commands, cmd)
				awaitingTitle = false
				awaitingDescription = false
				awaitingCommand = false
			} else {
				// Command is on the next line
				awaitingCommand = true
			}
			continue
		}
		
		// Handle lines following markers
		if trimmedLine != "" {
			if awaitingTitle {
				response.Title = trimmedLine
				awaitingTitle = false
				continue
			}
			if awaitingDescription {
				response.Description = trimmedLine
				awaitingDescription = false
				continue
			}
			if awaitingCommand {
				response.Commands = append(response.Commands, trimmedLine)
				awaitingCommand = false
				continue
			}
		}

		// Check for SCRIPT marker
		if strings.HasPrefix(trimmedLine, "SCRIPT:") {
			passedScriptMarker = true
			continue
		}

		// Collect script lines
		if passedScriptMarker {
			if trimmedLine != "" {
				scriptLines = append(scriptLines, line)
			}
		}
	}

	// Add script lines to commands
	if len(scriptLines) > 0 {
		response.Commands = append(response.Commands, strings.Join(scriptLines, "\n"))
	}

	// If no commands, try extracting code blocks
	if len(response.Commands) == 0 {
		response.Commands = extractCodeBlocks(rawResponse)
	}

	// If still no commands, return the entire response
	if len(response.Commands) == 0 && rawResponse != "" {
		response.Commands = []string{rawResponse}
	}

	return response
}

// Attempt to extract code blocks from markdown
func extractCodeBlocks(text string) []string {
	var commands []string
	lines := strings.Split(text, "\n")

	var inCodeBlock bool
	var codeLines []string

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if inCodeBlock {
				// Reached the end of the code block
				if len(codeLines) > 0 {
					// Join code lines and add as a single command to separate code blocks
					commands = append(commands, strings.Join(codeLines, "\n"))
					codeLines = nil
				}
				inCodeBlock = false
			} else {
				// Start of code block
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			codeLines = append(codeLines, line)
		}
	}

	return commands
}

// Create a new provider instance
func NewProvider(providerName, apiKey, model, baseURL string) (Provider, error) {
	switch providerName {
	case "OpenAI":
		return NewOpenAIProvider(apiKey, model), nil
	case "Anthropic":
		return NewAnthropicProvider(apiKey, model), nil
	case "Google":
		return NewGoogleProvider(apiKey, model), nil
	case "xAI":
		return NewXAIProvider(apiKey, model), nil
	case "OpenAI-Compatible":
		return NewOpenAICompatibleProvider(apiKey, model, baseURL), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", providerName)
	}
}
