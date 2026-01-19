package commands

import (
	"fmt"
	"os"

	"github.com/CRTYPUBG/winux/internal/utils"
)

// Pwd implements the pwd command.
// Usage: pwd
func Pwd(args []string) int {
	// Check for help
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			printPwdHelp()
			return utils.ExitSuccess
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
		return utils.ExitFailure
	}

	fmt.Println(dir)
	return utils.ExitSuccess
}

func printPwdHelp() {
	fmt.Println(`Usage: pwd

Print the full filename of the current working directory.

Options:
  --help    display this help and exit`)
}
