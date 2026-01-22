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
	// v0.1.0
	core.Register("ls", commands.Ls)
	core.Register("cat", commands.Cat)
	core.Register("grep", commands.Grep)
	// v0.2.0
	core.Register("rm", commands.Rm)
	core.Register("mkdir", commands.Mkdir)
	core.Register("touch", commands.Touch)
	core.Register("pwd", commands.Pwd)
	core.Register("echo", commands.Echo)
	// v0.3.0
	core.Register("whoami", commands.Whoami)
	core.Register("uptime", commands.Uptime)
	core.Register("nano", commands.Nano)
}

func main() {
	exitCode := core.Dispatch()
	os.Exit(exitCode)
}
