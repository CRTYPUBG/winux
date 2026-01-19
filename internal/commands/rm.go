package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/CRTYPUBG/winux/internal/utils"
)

// Rm implements the rm command.
// Usage: rm [-r] [-f] [-v] file...
func Rm(args []string) int {
	// Parse flags
	recursive := false // -r: recursive
	force := false     // -f: force (ignore errors)
	verbose := false   // -v: verbose

	var files []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-' {
			for _, ch := range arg[1:] {
				switch ch {
				case 'r', 'R':
					recursive = true
				case 'f':
					force = true
				case 'v':
					verbose = true
				default:
					fmt.Fprintf(os.Stderr, "rm: invalid option -- '%c'\n", ch)
					return utils.ExitUsageError
				}
			}
		} else if arg == "--recursive" {
			recursive = true
		} else if arg == "--force" {
			force = true
		} else if arg == "--verbose" {
			verbose = true
		} else if arg == "--help" {
			printRmHelp()
			return utils.ExitSuccess
		} else {
			files = append(files, arg)
		}
	}

	if len(files) == 0 {
		if !force {
			fmt.Fprintln(os.Stderr, "rm: missing operand")
			return utils.ExitUsageError
		}
		return utils.ExitSuccess
	}

	exitCode := utils.ExitSuccess

	for _, file := range files {
		info, err := os.Lstat(file)
		if err != nil {
			if !force {
				fmt.Fprintf(os.Stderr, "rm: cannot remove '%s': %v\n", file, err)
				exitCode = utils.ExitFailure
			}
			continue
		}

		if info.IsDir() {
			if !recursive {
				fmt.Fprintf(os.Stderr, "rm: cannot remove '%s': Is a directory\n", file)
				exitCode = utils.ExitFailure
				continue
			}

			err = os.RemoveAll(file)
		} else {
			err = os.Remove(file)
		}

		if err != nil {
			if !force {
				fmt.Fprintf(os.Stderr, "rm: cannot remove '%s': %v\n", file, err)
				exitCode = utils.ExitFailure
			}
		} else if verbose {
			fmt.Printf("removed '%s'\n", file)
		}
	}

	return exitCode
}

func printRmHelp() {
	fmt.Println(`Usage: rm [OPTION]... FILE...

Remove (unlink) the FILE(s).

Options:
  -f, --force      ignore nonexistent files, never prompt
  -r, -R, --recursive   remove directories and their contents recursively
  -v, --verbose    explain what is being done
  --help           display this help and exit

Examples:
  rm file.txt
  rm -f file.txt
  rm -rf directory/
  rm -v file1.txt file2.txt`)
}
