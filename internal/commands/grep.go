package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	winuxio "github.com/CRTYPUBG/winux/internal/io"
	"github.com/CRTYPUBG/winux/internal/utils"
)

// Grep implements the grep command.
// Usage: grep [-i] [-v] [-n] [-c] [-l] [-E] pattern [file...]
func Grep(args []string) int {
	// Parse flags
	ignoreCase := false    // -i: case insensitive
	invertMatch := false   // -v: invert match
	showLineNum := false   // -n: show line numbers
	countOnly := false     // -c: count matches only
	filesOnly := false     // -l: show filenames only
	extendedRegex := false // -E: extended regex (always enabled in Go)

	var pattern string
	var files []string

	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-' {
			for _, ch := range arg[1:] {
				switch ch {
				case 'i':
					ignoreCase = true
				case 'v':
					invertMatch = true
				case 'n':
					showLineNum = true
				case 'c':
					countOnly = true
				case 'l':
					filesOnly = true
				case 'E':
					extendedRegex = true
				case 'e':
					// -e pattern (next arg is pattern)
					i++
					if i < len(args) {
						pattern = args[i]
					}
				default:
					fmt.Fprintf(os.Stderr, "grep: invalid option -- '%c'\n", ch)
					return utils.ExitUsageError
				}
			}
		} else if arg == "--ignore-case" {
			ignoreCase = true
		} else if arg == "--invert-match" {
			invertMatch = true
		} else if arg == "--line-number" {
			showLineNum = true
		} else if arg == "--count" {
			countOnly = true
		} else if arg == "--files-with-matches" {
			filesOnly = true
		} else if arg == "--extended-regexp" {
			extendedRegex = true
		} else if arg == "--help" {
			printGrepHelp()
			return utils.ExitSuccess
		} else if pattern == "" {
			pattern = arg
		} else {
			files = append(files, arg)
		}
		i++
	}

	_ = extendedRegex // Go regex is always extended

	if pattern == "" {
		fmt.Fprintln(os.Stderr, "grep: missing pattern")
		printGrepHelp()
		return utils.ExitUsageError
	}

	// Build regex
	regexPattern := pattern
	if ignoreCase {
		regexPattern = "(?i)" + regexPattern
	}

	re, err := regexp.Compile(regexPattern)
	if err != nil {
		// Fall back to literal string matching
		re = nil
	}

	// If no files, read from stdin
	if len(files) == 0 {
		if winuxio.IsPiped() {
			files = []string{"-"}
		} else {
			fmt.Fprintln(os.Stderr, "grep: no input files")
			return utils.ExitUsageError
		}
	}

	matchFound := false
	multipleFiles := len(files) > 1

	for _, file := range files {
		var reader io.Reader
		var fileName string

		if file == "-" {
			reader = os.Stdin
			fileName = "(standard input)"
		} else {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", file, err)
				continue
			}
			defer f.Close()
			reader = f
			fileName = file
		}

		found := grepReader(reader, fileName, re, pattern, ignoreCase, invertMatch, showLineNum, countOnly, filesOnly, multipleFiles)
		if found {
			matchFound = true
		}
	}

	if matchFound {
		return utils.ExitSuccess
	}
	return utils.ExitFailure // No matches found
}

func grepReader(reader io.Reader, fileName string, re *regexp.Regexp, pattern string, ignoreCase, invertMatch, showLineNum, countOnly, filesOnly, showFileName bool) bool {
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineNum := 0
	matchCount := 0
	fileHasMatch := false

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		var match bool
		if re != nil {
			match = re.MatchString(line)
		} else {
			// Literal string match
			if ignoreCase {
				match = strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
			} else {
				match = strings.Contains(line, pattern)
			}
		}

		if invertMatch {
			match = !match
		}

		if match {
			matchCount++
			fileHasMatch = true

			if filesOnly {
				fmt.Println(fileName)
				return true // Only print filename once
			}

			if !countOnly {
				prefix := ""
				if showFileName {
					prefix = fileName + ":"
				}
				if showLineNum {
					prefix += fmt.Sprintf("%d:", lineNum)
				}
				fmt.Printf("%s%s\n", prefix, line)
			}
		}
	}

	if countOnly {
		if showFileName {
			fmt.Printf("%s:%d\n", fileName, matchCount)
		} else {
			fmt.Printf("%d\n", matchCount)
		}
	}

	return fileHasMatch
}

func printGrepHelp() {
	fmt.Println(`Usage: grep [OPTION]... PATTERN [FILE]...

Search for PATTERN in each FILE.
When FILE is -, read standard input.

Options:
  -i, --ignore-case         ignore case distinctions
  -v, --invert-match        select non-matching lines
  -n, --line-number         print line number with output lines
  -c, --count               print only a count of matching lines
  -l, --files-with-matches  print only names of files with matches
  -E, --extended-regexp     interpret pattern as extended regex
  -e PATTERN                use PATTERN for matching
  --help                    display this help and exit

Exit status:
  0  if any matches found
  1  if no matches found
  2  if error occurred

Examples:
  grep error log.txt
  grep -i ERROR log.txt
  grep -n "pattern" file1.txt file2.txt
  type log.txt | winux grep -i error`)
}
