package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/CRTYPUBG/winux/internal/utils"
)

// Touch implements the touch command.
// Usage: touch [-c] file...
func Touch(args []string) int {
	// Parse flags
	noCreate := false // -c: do not create new files

	var files []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-' {
			for _, ch := range arg[1:] {
				switch ch {
				case 'c':
					noCreate = true
				default:
					fmt.Fprintf(os.Stderr, "touch: invalid option -- '%c'\n", ch)
					return utils.ExitUsageError
				}
			}
		} else if arg == "--no-create" {
			noCreate = true
		} else if arg == "--help" {
			printTouchHelp()
			return utils.ExitSuccess
		} else {
			files = append(files, arg)
		}
	}

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "touch: missing file operand")
		return utils.ExitUsageError
	}

	exitCode := utils.ExitSuccess
	now := time.Now()

	for _, file := range files {
		_, err := os.Stat(file)
		fileExists := err == nil

		if !fileExists {
			if noCreate {
				continue
			}
			// Create new empty file
			f, err := os.Create(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "touch: cannot touch '%s': %v\n", file, err)
				exitCode = utils.ExitFailure
				continue
			}
			f.Close()
		} else {
			// Update timestamps
			err := os.Chtimes(file, now, now)
			if err != nil {
				fmt.Fprintf(os.Stderr, "touch: cannot touch '%s': %v\n", file, err)
				exitCode = utils.ExitFailure
			}
		}
	}

	return exitCode
}

func printTouchHelp() {
	fmt.Println(`Usage: touch [OPTION]... FILE...

Update the access and modification times of each FILE to the current time.
A FILE argument that does not exist is created empty.

Options:
  -c, --no-create  do not create any files
  --help           display this help and exit

Examples:
  touch file.txt
  touch -c existing.txt
  touch file1.txt file2.txt`)
}
