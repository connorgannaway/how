package ai

import (
	"fmt"
	"strings"

	"github.com/connorgannaway/how/internal/system"
)

// Creates a system prompt with OS and shell context
func BuildSystemPrompt(sysInfo *system.SystemInfo) string {
	return fmt.Sprintf(`You are a helpful terminal command assistant. The user is running:
OS: %s
Shell: %s
Package Manager: %s

When asked a question about how to do something in a terminal, respond with commands
and solutions SPECIFICALLY for this operating system and shell.

Important OS-specific considerations:
- macOS: Use homebrew, launchctl, system commands specific to Darwin
- Linux: Vary by distro (apt/yum/pacman), use systemctl, GNU coreutils
- Arch Linux: Use pacman, prefer Arch-specific approaches
- Ubuntu/Debian: Use apt, systemd
- Windows: Use PowerShell or cmd syntax, Windows-specific commands
- FreeBSD: Use pkg, rc.d, BSD-specific commands

Format your response as:
TITLE: [optional one-line title]
DESCRIPTION: [optional additional information]
COMMAND: [one-line command]

For multi-line scripts, use:
TITLE: [optional one-line title]
DESCRIPTION: [optional additional information]
SCRIPT:
[line 1]
[line 2]
[line 3]

Be concise and practical. Only include commands that directly answer the question
for the user's specific OS and shell. Only include a description for complex commands or scripts.
Do NOT include explanations outside the structured format. Keep it clean and executable.`,
sysInfo.OSName, sysInfo.Shell, sysInfo.GetPackageManager())
}

// Creates a user prompt from a question
func BuildUserPrompt(question string) string {
	if strings.HasPrefix(question, "how ") {
		return question
	}
	return "how " + question
}
