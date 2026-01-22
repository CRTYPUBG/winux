package commands

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/CRTYPUBG/winux/internal/utils"
)

var (
	modkernel32       = syscall.NewLazyDLL("kernel32.dll")
	procGetTickCount64 = modkernel32.NewProc("GetTickCount64")
)

// Uptime implements the uptime command for Windows.
func Uptime(args []string) int {
	// Check for help
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			printUptimeHelp()
			return utils.ExitSuccess
		}
	}

	ret, _, _ := procGetTickCount64.Call()
	if ret == 0 {
		fmt.Fprintf(os.Stderr, "uptime: failed to get system uptime\n")
		return utils.ExitFailure
	}

	uptimeDuration := time.Duration(ret) * time.Millisecond
	
	// Format: up 1 day, 2 hours, 30 minutes
	days := int(uptimeDuration.Hours()) / 24
	hours := int(uptimeDuration.Hours()) % 24
	minutes := int(uptimeDuration.Minutes()) % 60

	// Get current time
	now := time.Now().Format("15:04:05")

	fmt.Printf(" %s up ", now)
	if days > 0 {
		fmt.Printf("%d day(s), ", days)
	}
	fmt.Printf("%02d:%02d\n", hours, minutes)

	return utils.ExitSuccess
}

func printUptimeHelp() {
	fmt.Println(`Usage: uptime [OPTION]...
Display how long the system has been running.

Options:
  --help     display this help and exit`)
}
