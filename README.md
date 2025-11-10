# how

An AI-powered terminal command assistant for when [TLDR](https://github.com/tldr-pages/tldr) isn't helpful and the man page is too long.

![Demo](demo.gif)

## Installation

```bash
# go install
go install github.com/connorgannaway/how@latest

# homebrew
brew tap connorgannaway/tap
brew install how

# winget
winget source add -n connorgannaway https://github.com/connorgannaway/winget-pkgs
winget install connorgannaway.how

# scoop
scoop bucket add connorgannaway https://github.com/connorgannaway/scoops
scoop install how
```

## Quick Start

First time setup:

```bash
how --configure
```

Then start asking questions:

```bash
how do I check if a process is listening on port 3000
how do I get open pr authors using gh and jq
how do I find all files modified in the last 7 days
```

## Supported Providers

### Cloud Providers

Preconfigured cloud providers and models include

- **OpenAI** - GPT-5, GPT-4.1, GPT-4o, GPT-4, O-series models
- **Anthropic** - Claude Opus 4, Claude Sonnet 4.5, Claude 3.7/3.5
- **Google** - Gemini 2.5 Pro/Flash, Gemini 2.0 Flash
- **xAI** - Grok 4, Grok 3, Grok Code

The configuration flow includes a custom model field for models not listed.

### OpenAI-Compatible

Connect to any OpenAI-compatible API endpoint, including:

**Local:**

- [Ollama](https://ollama.ai) - `http://localhost:11434/v1`
- [LM Studio](https://lmstudio.ai) - `http://localhost:1234/v1`

**Cloud:**

- Groq
- DeepSeek
- Perplexity

## Usage

### Basic Usage

```bash
# Ask a question
how [question]

# Examples
how do I kill the process on port 8080
how do I compress png images over 20MB in a folder
```

### Configuration

![Configuration](configure.gif)

```bash
# Interactive configuration
how --configure
how -c

# View current configuration
how --status
how -s

# View with API key
how --status --key
how -s -k

# View all API keys
how --status --key --all
how -s -k -a

# Show full unmasked API keys
how --status --key --reveal-full
```

### Managing API Keys

```bash
# Interactive checklist to clear specific API keys
how --clear
how -r

# Clear all API keys at once
how --clear --all
how -r -a
```

### Other Commands

```bash
# Show version
how --version
how -v

# Show help
how --help
how -h
```

## Configuration

### Configuration File

Config file location defaults to XDG_CONFIG_HOME if set, otherwise:

- **macOS/Linux**: `~/.config/how/config.json`
- **Windows**: `%APPDATA%\how\config.json`

Example config:

```json
{
  "current_provider": "OpenAI-Compatible",
  "current_model": "gemma3:4b",
  "base_url": "http://localhost:11434/v1"
}
```

### API Key Storage

API keys are stored in the user keyring:

- **macOS**: Keychain
- **Windows**: Credential Manager
- **Linux**: Secret Service (gnome-keyring)

## License

MIT
