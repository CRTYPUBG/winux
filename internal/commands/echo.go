package commands

import (
	"fmt"
	"strings"

	"github.com/CRTYPUBG/winux/internal/utils"
)

// Echo implements the echo command.
// Usage: echo [-n] [-e] [string...]
func Echo(args []string) int {
	// Parse flags
	noNewline := false      // -n: no trailing newline
	interpretEscapes := false // -e: interpret escape sequences

	var parts []string
	flagsDone := false

	for _, arg := range args {
		if !flagsDone && strings.HasPrefix(arg, "-") && len(arg) > 1 {
			allFlags := true
			for _, ch := range arg[1:] {
				switch ch {
				case 'n':
					noNewline = true
				case 'e':
					interpretEscapes = true
				case 'E':
					interpretEscapes = false
				case '-':
					// -- ends flag processing
					flagsDone = true
				default:
					allFlags = false
				}
			}
			if allFlags && arg != "--" {
				continue
			}
			if arg == "--" {
				flagsDone = true
				continue
			}
		}
		flagsDone = true
		parts = append(parts, arg)
	}

	output := strings.Join(parts, " ")

	if interpretEscapes {
		output = interpretEscapeSequences(output)
	}

	if noNewline {
		fmt.Print(output)
	} else {
		fmt.Println(output)
	}

	return utils.ExitSuccess
}

func interpretEscapeSequences(s string) string {
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result.WriteByte('\n')
				i += 2
			case 't':
				result.WriteByte('\t')
				i += 2
			case 'r':
				result.WriteByte('\r')
				i += 2
			case '\\':
				result.WriteByte('\\')
				i += 2
			case '0':
				result.WriteByte(0)
				i += 2
			case 'a':
				result.WriteByte('\a')
				i += 2
			case 'b':
				result.WriteByte('\b')
				i += 2
			case 'f':
				result.WriteByte('\f')
				i += 2
			case 'v':
				result.WriteByte('\v')
				i += 2
			default:
				result.WriteByte(s[i])
				i++
			}
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}

func printEchoHelp() {
	fmt.Println(`Usage: echo [OPTION]... [STRING]...

Echo the STRING(s) to standard output.

Options:
  -n    do not output the trailing newline
  -e    enable interpretation of backslash escapes
  -E    disable interpretation of backslash escapes (default)

Escape sequences (with -e):
  \\    backslash
  \n    new line
  \t    horizontal tab
  \r    carriage return

Examples:
  echo Hello World
  echo -n "no newline"
  echo -e "line1\nline2"`)
}
