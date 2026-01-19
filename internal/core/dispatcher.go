// Package core provides the command dispatcher and registry.
package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CRTYPUBG/winux/internal/utils"
)

// CommandFunc is the signature for all command implementations.
// Returns an exit code.
type CommandFunc func(args []string) int

// Registry holds all registered commands.
var Registry = make(map[string]CommandFunc)

// Register adds a command to the registry.
func Register(name string, fn CommandFunc) {
	Registry[name] = fn
}

// Dispatch resolves and executes the appropriate command.
// Resolution order:
//  1. Executable name (argv[0]) - BusyBox style
//  2. First CLI argument (winux <command>)
func Dispatch() int {
	// Try argv[0] first (BusyBox-style symlink dispatch)
	execName := filepath.Base(os.Args[0])
	execName = strings.TrimSuffix(execName, ".exe")
	execName = strings.ToLower(execName)

	// If invoked as a command directly (e.g., "grep.exe" or symlink "grep")
	if execName != "winux" {
		if fn, ok := Registry[execName]; ok {
			return fn(os.Args[1:])
		}
	}

	// Otherwise, expect "winux <command> [args...]"
	if len(os.Args) < 2 {
		printUsage()
		return utils.ExitUsageError
	}

	cmdName := strings.ToLower(os.Args[1])

	// Handle help flags
	if cmdName == "--help" || cmdName == "-h" || cmdName == "help" {
		printUsage()
		return utils.ExitSuccess
	}

	// Handle version
	if cmdName == "--version" || cmdName == "-v" || cmdName == "version" {
		fmt.Println("winux v0.1.0")
		return utils.ExitSuccess
	}

	// Look up command
	if fn, ok := Registry[cmdName]; ok {
		return fn(os.Args[2:])
	}

	fmt.Fprintf(os.Stderr, "winux: '%s' is not a winux command. See 'winux --help'.\n", cmdName)
	return utils.ExitCommandNotFound
}

func printUsage() {
	fmt.Println(`WINUX - Native Linux-like utilities for Windows

Usage: winux <command> [arguments]

Available commands:
  cat      Concatenate and print files
  echo     Display a line of text
  grep     Search for patterns in files
  ls       List directory contents
  mkdir    Create directories
  pwd      Print working directory
  rm       Remove files or directories
  touch    Create files or update timestamps

Options:
  --help, -h       Show this help message
  --version, -v    Show version information

Examples:
  winux ls -la
  winux cat file.txt
  winux grep -i error log.txt
  winux mkdir -p path/to/dir
  winux rm -rf temp/
  type log.txt | winux grep error`)
}
