package clipboard

import (
	"strings"

	"github.com/atotto/clipboard"
)

// clipboard wrapper for copying string arrays to clipboard
func CopyCommands(commands []string) error {
	text := strings.Join(commands, "\n")
	return clipboard.WriteAll(text)
}
