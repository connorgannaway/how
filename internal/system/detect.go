package system

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type SystemInfo struct {
	OS        string // "darwin", "linux", "windows", "freebsd", etc.
	OSName    string // "macOS", "Ubuntu", "Windows", "FreeBSD", etc.
	Shell     string // "bash", "zsh", "fish", "powershell", "cmd", etc.
	ShellPath string // Full path to shell executable
}

// Detect the current operating system and shell
func DetectSystem() (*SystemInfo, error) {
	info := &SystemInfo{
		OS: runtime.GOOS,
	}

	// Detect OS name
	switch runtime.GOOS {
	case "darwin":
		info.OSName = "macOS"
	case "linux":
		info.OSName = detectLinuxDistro()
	case "windows":
		info.OSName = "Windows"
	case "freebsd":
		info.OSName = "FreeBSD"
	default:
		info.OSName = runtime.GOOS
	}

	// Detect shell
	info.Shell, info.ShellPath = detectShell()

	return info, nil
}

// Returns a formatted string for AI context
func (s *SystemInfo) GetContextString() string {
	return fmt.Sprintf("%s (%s) using %s shell", s.OSName, s.OS, s.Shell)
}

// Attempt to detect the Linux distribution
func detectLinuxDistro() string {
	// Try to read /etc/os-release
	data, err := os.ReadFile("/etc/os-release")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "NAME=") {
				name := strings.TrimPrefix(line, "NAME=")
				name = strings.Trim(name, "\"")
				return name
			}
		}
	}

	// Try /etc/lsb-release
	data, err = os.ReadFile("/etc/lsb-release")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "DISTRIB_ID=") {
				name := strings.TrimPrefix(line, "DISTRIB_ID=")
				name = strings.Trim(name, "\"")
				return name
			}
		}
	}

	return "Linux"
}

// Attempt to detect the current shell
func detectShell() (string, string) {
	// Check SHELL environment variable
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		// On Windows, check for PowerShell or cmd
		if runtime.GOOS == "windows" {
			if os.Getenv("PSModulePath") != "" {
				return "powershell", "powershell.exe"
			}
			return "cmd", "cmd.exe"
		}
		// Default to sh on Unix-like systems
		return "sh", "/bin/sh"
	}

	// If we have shellpath, extract shell name from path
	parts := strings.Split(shellPath, "/")
	shellName := parts[len(parts)-1]

	return shellName, shellPath
}

// Return the likely package manager for the system
func (s *SystemInfo) GetPackageManager() string {
	switch s.OS {
	case "darwin":
		return "brew"
	case "windows":
		if commandExists("choco") {
			return "choco"
		}
		if commandExists("scoop") {
			return "scoop"
		}
		// Default to winget 
		return "winget"
	case "freebsd":
		return "pkg"
	case "linux":
		if commandExists("apt") {
			return "apt"
		}
		if commandExists("pacman") {
			return "pacman"
		}
		if commandExists("yum") {
			return "yum"
		}
		if commandExists("dnf") {
			return "dnf"
		}
		if commandExists("zypper") {
			return "zypper"
		}
		return "unknown"
	default:
		return "unknown"
	}
}

// Check if a command exists in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
