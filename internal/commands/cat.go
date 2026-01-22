package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	winuxio "github.com/CRTYPUBG/winux/internal/io"
	"github.com/CRTYPUBG/winux/internal/utils"
)

// Cat implements the cat command.
// Usage: cat [-n] [-b] [file...]
func Cat(args []string) int {
	// Parse flags
	numberLines := false    // -n: number all output lines
	numberNonBlank := false // -b: number non-blank output lines

	var files []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-' {
			for _, ch := range arg[1:] {
				switch ch {
				case 'n':
					numberLines = true
				case 'b':
					numberNonBlank = true
				default:
					fmt.Fprintf(os.Stderr, "cat: invalid option -- '%c'\n", ch)
					return utils.ExitUsageError
				}
			}
		} else if arg == "--number" {
			numberLines = true
		} else if arg == "--number-nonblank" {
			numberNonBlank = true
		} else if arg == "--help" {
			printCatHelp()
			return utils.ExitSuccess
		} else if arg == "-" {
			files = append(files, "-") // stdin marker
		} else {
			files = append(files, arg)
		}
	}

	// -b overrides -n
	if numberNonBlank {
		numberLines = false
	}

	// If no files specified, read from stdin
	if len(files) == 0 {
		if winuxio.IsPiped() {
			files = []string{"-"}
		} else {
			// No input at all
			printCatHelp()
			return utils.ExitUsageError
		}
	}

	exitCode := utils.ExitSuccess
	lineNum := 1

	for _, file := range files {
		var r io.Reader
		var closer io.Closer

		if file == "-" {
			r = os.Stdin
		} else {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cat: %s: %v\n", file, err)
				exitCode = utils.ExitFailure
				continue
			}
			r = f
			closer = f
		}

		if !numberLines && !numberNonBlank {
			// Fast path for raw output (supports binary files)
			_, err := io.Copy(os.Stdout, r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cat: %s: %v\n", file, err)
				exitCode = utils.ExitFailure
			}
		} else {
			// Line numbering path
			scanner := bufio.NewScanner(r)
			buf := make([]byte, 0, 64*1024)
			scanner.Buffer(buf, 1024*1024)

			for scanner.Scan() {
				line := scanner.Text()
				if numberNonBlank {
					if line != "" {
						fmt.Printf("%6d\t%s\n", lineNum, line)
						lineNum++
					} else {
						fmt.Println()
					}
				} else if numberLines {
					fmt.Printf("%6d\t%s\n", lineNum, line)
					lineNum++
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "cat: %s: %v\n", file, err)
				exitCode = utils.ExitFailure
			}
		}

		if closer != nil {
			closer.Close()
		}
	}

	return exitCode
}

func printCatHelp() {
	fmt.Println(`Usage: cat [OPTION]... [FILE]...

Concatenate FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

Options:
  -b, --number-nonblank    number nonempty output lines
  -n, --number             number all output lines
  --help                   display this help and exit

Examples:
  cat file.txt
  cat -n file.txt
  type file.txt | winux cat -n`)
}
