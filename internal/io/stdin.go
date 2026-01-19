// Package io provides IO and pipe handling utilities.
package io

import (
	"os"
)

// IsPiped returns true if stdin is receiving piped input.
// This enables: type file.txt | winux grep error
func IsPiped() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}
