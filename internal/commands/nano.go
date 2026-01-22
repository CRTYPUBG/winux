package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/CRTYPUBG/winux/internal/utils"
)

var (
	kernel32                      = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode            = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode            = kernel32.NewProc("SetConsoleMode")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

type consoleScreenBufferInfo struct {
	Size              coord
	CursorPosition    coord
	Attributes        uint16
	Window            smallRect
	MaximumWindowSize coord
}

type coord struct {
	X, Y int16
}

type smallRect struct {
	Left, Top, Right, Bottom int16
}

const (
	enableLineInput       = 0x0002
	enableEchoInput       = 0x0004
	enableProcessedInput  = 0x0001
	enableExtendedFlags   = 0x0080
	enableVirtualTerminal = 0x0004 // Enable ANSI escape sequences on output
)

type Editor struct {
	lines       []string
	cursorX     int
	cursorY     int
	offsetX     int
	offsetY     int
	width       int
	height      int
	filename    string
	dirty       bool
	statusMsg   string
	originalIn  uint32
	originalOut uint32
}

func Nano(args []string) int {
	if len(args) < 1 {
		fmt.Println("Usage: nano <filename>")
		return utils.ExitUsageError
	}

	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			printNanoHelp()
			return utils.ExitSuccess
		}
	}

	filename := args[0]
	e := &Editor{filename: filename}
	e.load()

	if err := e.enterRawMode(); err != nil {
		fmt.Fprintf(os.Stderr, "nano: failed to enter raw mode: %v\n", err)
		return utils.ExitFailure
	}
	defer e.exitRawMode()

	// Main loop
	for {
		e.updateSize()
		e.refreshScreen()
		if !e.processInput() {
			break
		}
	}

	return utils.ExitSuccess
}

func (e *Editor) load() {
	data, err := os.ReadFile(e.filename)
	if err != nil {
		e.lines = []string{""}
		e.statusMsg = "New File: " + e.filename
		return
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		lines[i] = strings.ReplaceAll(line, "\r", "")
	}
	e.lines = lines
	e.statusMsg = fmt.Sprintf("Read %d lines", len(lines))
}

func (e *Editor) save() {
	content := strings.Join(e.lines, "\n")
	err := os.WriteFile(e.filename, []byte(content), 0644)
	if err != nil {
		e.statusMsg = "Error saving: " + err.Error()
	} else {
		e.dirty = false
		e.statusMsg = "Saved " + e.filename
	}
}

func (e *Editor) enterRawMode() error {
	in := syscall.Handle(os.Stdin.Fd())
	out := syscall.Handle(os.Stdout.Fd())

	procGetConsoleMode.Call(uintptr(in), uintptr(unsafe.Pointer(&e.originalIn)))
	procGetConsoleMode.Call(uintptr(out), uintptr(unsafe.Pointer(&e.originalOut)))

	// Disable echo and line input, enable ANSI processing
	newIn := e.originalIn &^ (enableLineInput | enableEchoInput | enableProcessedInput)
	procSetConsoleMode.Call(uintptr(in), uintptr(newIn))

	newOut := e.originalOut | enableVirtualTerminal
	procSetConsoleMode.Call(uintptr(out), uintptr(newOut))

	// Clear screen and enter alternative buffer (simulated)
	fmt.Print("\033[?1049h") // Alternate buffer
	return nil
}

func (e *Editor) exitRawMode() {
	fmt.Print("\033[?1049l") // Main buffer
	in := syscall.Handle(os.Stdin.Fd())
	out := syscall.Handle(os.Stdout.Fd())
	procSetConsoleMode.Call(uintptr(in), uintptr(e.originalIn))
	procSetConsoleMode.Call(uintptr(out), uintptr(e.originalOut))
}

func (e *Editor) updateSize() {
	var info consoleScreenBufferInfo
	out := syscall.Handle(os.Stdout.Fd())
	procGetConsoleScreenBufferInfo.Call(uintptr(out), uintptr(unsafe.Pointer(&info)))
	e.width = int(info.Window.Right - info.Window.Left + 1)
	e.height = int(info.Window.Bottom - info.Window.Top + 1)
}

func (e *Editor) refreshScreen() {
	var sb strings.Builder
	sb.WriteString("\033[H") // Move cursor to top-left

	// 1. Header
	header := fmt.Sprintf(" WINUX nano %s ", e.filename)
	if e.dirty {
		header += "(modified) "
	}
	padding := e.width - len(header)
	if padding < 0 {
		padding = 0
	}
	sb.WriteString("\033[7m") // Invert
	sb.WriteString(header)
	sb.WriteString(strings.Repeat(" ", padding))
	sb.WriteString("\033[0m\r\n")

	// 2. Content
	viewHeight := e.height - 4 // Header(1) + Status(1) + Help(2)
	for i := 0; i < viewHeight; i++ {
		lineIdx := i + e.offsetY
		if lineIdx < len(e.lines) {
			line := e.lines[lineIdx]
			if len(line) > e.offsetX {
				line = line[e.offsetX:]
			} else {
				line = ""
			}
			if len(line) > e.width {
				line = line[:e.width]
			}
			sb.WriteString(line)
		}
		sb.WriteString("\033[K") // Clear to end of line
		sb.WriteString("\r\n")
	}

	// 3. Status Bar
	sb.WriteString("\033[7m")
	status := fmt.Sprintf(" Line: %d/%d Col: %d ", e.cursorY+1, len(e.lines), e.cursorX+1)
	sb.WriteString(status)
	sb.WriteString(strings.Repeat(" ", e.width-len(status)))
	sb.WriteString("\033[0m\r\n")

	// 4. Help Message
	sb.WriteString("\033[K") // Clear line
	sb.WriteString(e.statusMsg)
	sb.WriteString("\r\n")
	sb.WriteString("\033[7m^X\033[0m Exit  \033[7m^O\033[0m Save")

	// 5. Position Cursor
	sb.WriteString(fmt.Sprintf("\033[%d;%dH", e.cursorY-e.offsetY+2, e.cursorX-e.offsetX+1))

	fmt.Print(sb.String())
}

func (e *Editor) processInput() bool {
	reader := bufio.NewReader(os.Stdin)
	b, err := reader.ReadByte()
	if err != nil {
		return false
	}

	switch b {
	case 24: // Ctrl+X
		if e.dirty {
			// In a real nano, this would ask to save. For minimalism, we just exit.
			// But let's stay friendly.
		}
		return false
	case 15: // Ctrl+O
		e.save()
	case 13, 10: // Enter
		currLine := e.lines[e.cursorY]
		nextLine := currLine[e.cursorX:]
		e.lines[e.cursorY] = currLine[:e.cursorX]
		
		newLines := make([]string, 0, len(e.lines)+1)
		newLines = append(newLines, e.lines[:e.cursorY+1]...)
		newLines = append(newLines, nextLine)
		newLines = append(newLines, e.lines[e.cursorY+1:]...)
		e.lines = newLines
		
		e.cursorY++
		e.cursorX = 0
		e.dirty = true
	case 8, 127: // Backspace
		if e.cursorX > 0 {
			line := e.lines[e.cursorY]
			e.lines[e.cursorY] = line[:e.cursorX-1] + line[e.cursorX:]
			e.cursorX--
			e.dirty = true
		} else if e.cursorY > 0 {
			prevLineLen := len(e.lines[e.cursorY-1])
			e.lines[e.cursorY-1] += e.lines[e.cursorY]
			e.lines = append(e.lines[:e.cursorY], e.lines[e.cursorY+1:]...)
			e.cursorY--
			e.cursorX = prevLineLen
			e.dirty = true
		}
	case 27: // Escape / Arrows
		// Simplified arrow handling
		nextByte, _ := reader.ReadByte()
		if nextByte == '[' {
			arrow, _ := reader.ReadByte()
			switch arrow {
			case 'A': // Up
				if e.cursorY > 0 {
					e.cursorY--
				}
			case 'B': // Down
				if e.cursorY < len(e.lines)-1 {
					e.cursorY++
				}
			case 'C': // Right
				if e.cursorX < len(e.lines[e.cursorY]) {
					e.cursorX++
				}
			case 'D': // Left
				if e.cursorX > 0 {
					e.cursorX--
				}
			}
			// Keep cursor valid
			if e.cursorX > len(e.lines[e.cursorY]) {
				e.cursorX = len(e.lines[e.cursorY])
			}
		}
	default:
		// Regular char
		if b >= 32 && b <= 126 {
			line := e.lines[e.cursorY]
			e.lines[e.cursorY] = line[:e.cursorX] + string(b) + line[e.cursorX:]
			e.cursorX++
			e.dirty = true
		}
	}

	// Scroll management
	viewHeight := e.height - 4
	if e.cursorY < e.offsetY {
		e.offsetY = e.cursorY
	}
	if e.cursorY >= e.offsetY+viewHeight {
		e.offsetY = e.cursorY - viewHeight + 1
	}
	if e.cursorX < e.offsetX {
		e.offsetX = e.cursorX
	}
	if e.cursorX >= e.offsetX+e.width {
		e.offsetX = e.cursorX - e.width + 1
	}

	return true
}

func printNanoHelp() {
	fmt.Println(`Usage: nano [FILE]
A minimal terminal text editor for Windows.

Options:
  --help     display this help and exit

Keybindings:
  ^X         Exit
  ^O         Save
  Arrows     Navigate`)
}
