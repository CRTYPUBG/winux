package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/CRTYPUBG/winux/internal/utils"
)

// Mkdir implements the mkdir command.
// Usage: mkdir [-p] [-v] directory...
func Mkdir(args []string) int {
	// Parse flags
	parents := false // -p: create parent directories
	verbose := false // -v: verbose

	var dirs []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-' {
			for _, ch := range arg[1:] {
				switch ch {
				case 'p':
					parents = true
				case 'v':
					verbose = true
				default:
					fmt.Fprintf(os.Stderr, "mkdir: invalid option -- '%c'\n", ch)
					return utils.ExitUsageError
				}
			}
		} else if arg == "--parents" {
			parents = true
		} else if arg == "--verbose" {
			verbose = true
		} else if arg == "--help" {
			printMkdirHelp()
			return utils.ExitSuccess
		} else {
			dirs = append(dirs, arg)
		}
	}

	if len(dirs) == 0 {
		fmt.Fprintln(os.Stderr, "mkdir: missing operand")
		return utils.ExitUsageError
	}

	exitCode := utils.ExitSuccess

	for _, dir := range dirs {
		var err error

		if parents {
			err = os.MkdirAll(dir, 0755)
		} else {
			err = os.Mkdir(dir, 0755)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "mkdir: cannot create directory '%s': %v\n", dir, err)
			exitCode = utils.ExitFailure
		} else if verbose {
			fmt.Printf("mkdir: created directory '%s'\n", dir)
		}
	}

	return exitCode
}

func printMkdirHelp() {
	fmt.Println(`Usage: mkdir [OPTION]... DIRECTORY...

Create the DIRECTORY(ies), if they do not already exist.

Options:
  -p, --parents    no error if existing, make parent directories as needed
  -v, --verbose    print a message for each created directory
  --help           display this help and exit

Examples:
  mkdir newdir
  mkdir -p path/to/newdir
  mkdir -v dir1 dir2 dir3`)
}
