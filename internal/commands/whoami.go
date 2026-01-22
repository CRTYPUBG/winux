package commands

import (
	"fmt"
	"os"
	"os/user"

	"github.com/CRTYPUBG/winux/internal/utils"
)

// Whoami implements the whoami command.
// Usage: whoami
func Whoami(args []string) int {
	// Check for help
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			printWhoamiHelp()
			return utils.ExitSuccess
		}
	}

	currUser, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "whoami: %v\n", err)
		return utils.ExitFailure
	}

	// On Windows, Username might include domain (DOMAIN\User).
	// Linux whoami usually just returns the user part.
	username := currUser.Username
	if lastSlash := lastIndex(username, '\\'); lastSlash != -1 {
		username = username[lastSlash+1:]
	}

	fmt.Println(username)
	return utils.ExitSuccess
}

func lastIndex(s string, char byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == char {
			return i
		}
	}
	return -1
}

func printWhoamiHelp() {
	fmt.Println(`Usage: whoami [OPTION]...
Print the user name associated with the current effective user ID.

Options:
  --help     display this help and exit`)
}
