package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/connorgannaway/how/internal/ai"
	"github.com/connorgannaway/how/internal/config"
	"github.com/connorgannaway/how/internal/system"
	"github.com/connorgannaway/how/internal/ui/clear"
	"github.com/connorgannaway/how/internal/ui/configure"
	"github.com/connorgannaway/how/internal/ui/question"
	"github.com/connorgannaway/how/internal/ui/status"
)

var Version = "dev"

func main() {
	// Define flags
	versionFlag := flag.Bool("v", false, "Print version and exit")
	versionLongFlag := flag.Bool("version", false, "Print version and exit")
	configureFlag := flag.Bool("c", false, "Configure AI provider and API key")
	configureLongFlag := flag.Bool("configure", false, "Configure AI provider and API key")
	statusFlag := flag.Bool("s", false, "Show current configuration status")
	statusLongFlag := flag.Bool("status", false, "Show current configuration status")
	keyFlag := flag.Bool("k", false, "Show API key(s) with --status")
	keyLongFlag := flag.Bool("key", false, "Show API key(s) with --status")
	revealFullFlag := flag.Bool("reveal-full", false, "Show full unmasked API keys with --status --key")
	clearFlag := flag.Bool("r", false, "Clear API keys from configuration")
	clearLongFlag := flag.Bool("clear", false, "Clear API keys from configuration")
	allFlag := flag.Bool("a", false, "With --status --key: show all provider API keys. With --clear: clear all API keys without prompting")
	allLongFlag := flag.Bool("all", false, "With --status --key: show all provider API keys. With --clear: clear all API keys without prompting")
	helpFlag := flag.Bool("h", false, "Show help message")
	helpLongFlag := flag.Bool("help", false, "Show help message")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: how [options] <question>\n\n")
		fmt.Fprintf(os.Stderr, "AI-powered terminal command assistant\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -c, --configure    Configure AI provider and API key\n")
		fmt.Fprintf(os.Stderr, "  -s, --status       Show current configuration status\n")
		fmt.Fprintf(os.Stderr, "  -k, --key          Show API key(s) with --status (masked by default)\n")
		fmt.Fprintf(os.Stderr, "  --reveal-full      Show full unmasked API keys (use with --status --key)\n")
		fmt.Fprintf(os.Stderr, "  -r, --clear        Clear API keys from configuration\n")
		fmt.Fprintf(os.Stderr, "  -a, --all          Use with --status/--clear for all providers\n")
		fmt.Fprintf(os.Stderr, "  -v, --version      Print version and exit\n")
		fmt.Fprintf(os.Stderr, "  -h, --help         Show this help message\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  how do I check if a process is listening on port 3000\n")
		fmt.Fprintf(os.Stderr, "  how do I compress png images over 20MB in a folder\n")
		fmt.Fprintf(os.Stderr, "  how --configure\n")
		fmt.Fprintf(os.Stderr, "  how --status\n")
		fmt.Fprintf(os.Stderr, "  how --status --key\n")
		fmt.Fprintf(os.Stderr, "  how --status --key --reveal-full    # Shows full keys\n")
		fmt.Fprintf(os.Stderr, "  how --status --key --all\n")
		fmt.Fprintf(os.Stderr, "  how --clear\n")
		fmt.Fprintf(os.Stderr, "  how --clear --all\n")
	}

	flag.Parse()

	// Handle version flag
	if *versionFlag || *versionLongFlag {
		fmt.Printf("how version %s\n", Version)
		os.Exit(0)
	}

	// Handle help flag
	if *helpFlag || *helpLongFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Handle configure flag
	if *configureFlag || *configureLongFlag {
		if err := configure.Run(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error during configuration: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Handle clear flag
	if *clearFlag || *clearLongFlag {
		if *allFlag || *allLongFlag {
			// Clear all API keys without prompting
			providers := config.GetProviders()
			cleared := 0
			for _, provider := range providers {
				hasKey, _ := config.HasAPIKeyInKeyring(provider)
				if hasKey {
					if err := config.DeleteAPIKeyFromKeyring(provider); err != nil {
						fmt.Fprintf(os.Stderr, "Error deleting API key for %s: %v\n", provider, err)
						os.Exit(1)
					}
					cleared++
				}
			}

			if cleared > 0 {
				fmt.Printf("✓ Cleared %d API key(s)\n", cleared)
			} else {
				fmt.Println("No API keys to clear.")
			}
		} else {
			// Show interactive checklist
			if err := clear.Run(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error during clear: %v\n", err)
				os.Exit(1)
			}
		}
		os.Exit(0)
	}

	// Handle status flag
	if *statusFlag || *statusLongFlag {
		status.Run(cfg, *keyFlag || *keyLongFlag, *allFlag || *allLongFlag, *revealFullFlag)
		os.Exit(0)
	}

	// Check if configured
	if ready, missing := cfg.IsConfigured(); !ready {
		fmt.Fprintf(os.Stderr, "Not configured. Missing: %v. Run 'how --configure' to set up.\n", strings.Join(missing, ", "))
		os.Exit(1)
	}

	// Build question from arguments
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	questionText := strings.Join(args, " ")

	// Detect system
	sysInfo, err := system.DetectSystem()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting system: %v\n", err)
		os.Exit(1)
	}

	// Retrieve API key
	apiKey, err := config.GetAPIKeyFromKeyring(cfg.CurrentProvider)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving API key: %v\n", err)
		os.Exit(1)
	}

	// Create AI provider
	provider, err := ai.NewProvider(
		cfg.CurrentProvider,
		apiKey,
		cfg.CurrentModel,
		cfg.BaseURL,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating AI provider: %v\n", err)
		os.Exit(1)
	}

	// Run question UI
	if err := question.Run(questionText, provider, sysInfo); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
