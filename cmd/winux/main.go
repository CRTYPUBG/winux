// WINUX - Native Linux-like core utilities for Windows
// Single binary · Go · No WSL · No aliases
//
// Entry point and command registration.
// Protected binary - anti-debugging enabled.
package main

import (
	"os"
	"runtime"

	"github.com/CRTYPUBG/winux/internal/commands"
	"github.com/CRTYPUBG/winux/internal/core"
	"github.com/CRTYPUBG/winux/internal/protection"
)

// Build-time variables (set via ldflags)
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func init() {
	// =========================================
	// SECURITY: Anti-debug protection
	// =========================================
	// Initialize protection before anything else
	// This will terminate if debugger is detected
	protection.Init()

	// Use all available CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Register all commands
	core.Register("ls", commands.Ls)
	core.Register("cat", commands.Cat)
	core.Register("grep", commands.Grep)
}

func main() {
	exitCode := core.Dispatch()
	os.Exit(exitCode)
}
