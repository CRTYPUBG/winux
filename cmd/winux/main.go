// WINUX - Native Linux-like core utilities for Windows
// Single binary · Go · No WSL · No aliases
//
// Entry point and command registration.
package main

import (
	"os"
	"runtime"

	"github.com/CRTYPUBG/winux/internal/commands"
	"github.com/CRTYPUBG/winux/internal/core"
)

func init() {
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
