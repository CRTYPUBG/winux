// Package utils provides shared utilities for WINUX commands.
package utils

// Linux-compatible exit codes
const (
	ExitSuccess         = 0   // Command completed successfully
	ExitFailure         = 1   // General failure (e.g., no matches found)
	ExitUsageError      = 2   // Invalid usage, missing args, permission denied
	ExitCommandNotFound = 127 // Command not recognized
)
