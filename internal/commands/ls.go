// Package commands contains all WINUX command implementations.
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/CRTYPUBG/winux/internal/utils"
)

// Ls implements the ls command.
// Usage: ls [-l] [-a] [-h] [path...]
func Ls(args []string) int {
	// Parse flags
	showAll := false      // -a: show hidden files
	longFormat := false   // -l: long listing format
	humanReadable := false // -h: human readable sizes
	var paths []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-' {
			// Short flags
			for _, ch := range arg[1:] {
				switch ch {
				case 'a':
					showAll = true
				case 'l':
					longFormat = true
				case 'h':
					humanReadable = true
				default:
					fmt.Fprintf(os.Stderr, "ls: invalid option -- '%c'\n", ch)
					return utils.ExitUsageError
				}
			}
		} else if arg == "--all" {
			showAll = true
		} else if arg == "--help" {
			printLsHelp()
			return utils.ExitSuccess
		} else {
			paths = append(paths, arg)
		}
	}

	// Default to current directory
	if len(paths) == 0 {
		paths = []string{"."}
	}

	exitCode := utils.ExitSuccess
	multiplePaths := len(paths) > 1

	for i, path := range paths {
		if multiplePaths {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", path)
		}

		code := listDir(path, showAll, longFormat, humanReadable)
		if code != utils.ExitSuccess {
			exitCode = code
		}
	}

	return exitCode
}

func listDir(path string, showAll, longFormat, humanReadable bool) int {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %v\n", path, err)
		return utils.ExitUsageError
	}

	// Sort entries alphabetically (case-insensitive)
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files unless -a
		if !showAll && strings.HasPrefix(name, ".") {
			continue
		}

		if longFormat {
			info, err := entry.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "ls: cannot stat '%s': %v\n", filepath.Join(path, name), err)
				continue
			}

			// Format: permissions size date name
			mode := info.Mode().String()
			size := info.Size()
			modTime := info.ModTime().Format("Jan _2 15:04")

			sizeStr := fmt.Sprintf("%8d", size)
			if humanReadable {
				sizeStr = fmt.Sprintf("%8s", formatSize(size))
			}

			// Directory indicator
			if entry.IsDir() {
				name = name + "/"
			}

			fmt.Printf("%s %s %s %s\n", mode, sizeStr, modTime, name)
		} else {
			if entry.IsDir() {
				fmt.Printf("%s/\n", name)
			} else {
				fmt.Println(name)
			}
		}
	}

	return utils.ExitSuccess
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1fG", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1fM", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1fK", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

func printLsHelp() {
	fmt.Println(`Usage: ls [OPTION]... [FILE]...

List directory contents.

Options:
  -a, --all    do not ignore entries starting with .
  -l           use a long listing format
  -h           with -l, print sizes in human readable format
  --help       display this help and exit`)
}
