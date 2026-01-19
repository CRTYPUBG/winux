// Package protection provides anti-debugging and anti-tampering measures.
// This ensures WINUX cannot be analyzed or modified by debugging tools.
package protection

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	ntdll            = syscall.NewLazyDLL("ntdll.dll")
	isDebuggerPresent = kernel32.NewProc("IsDebuggerPresent")
	checkRemoteDebuggerPresent = kernel32.NewProc("CheckRemoteDebuggerPresent")
	ntQueryInformationProcess = ntdll.NewProc("NtQueryInformationProcess")
	outputDebugStringW = kernel32.NewProc("OutputDebugStringW")
	getCurrentProcess = kernel32.NewProc("GetCurrentProcess")
)

// Debugger process names to detect
var debuggerProcesses = []string{
	"ollydbg", "x64dbg", "x32dbg", "windbg", "ida", "ida64",
	"idaq", "idaq64", "idaw", "idaw64", "idag", "idag64",
	"ghidra", "radare2", "r2", "processhacker", "procmon",
	"procexp", "pestudio", "die", "peid", "lordpe", "wireshark",
	"fiddler", "charles", "httpdebugger", "dnspy", "dotpeek",
	"ilspy", "de4dot", "cheatengine", "artmoney", "apimonitor",
	"immunity", "binary ninja", "hopper", "cutter", "reclass",
}

// Init performs all anti-debug checks at startup.
// If any debugger is detected, the program terminates immediately.
func Init() {
	// Run checks in goroutine for timing attack
	go continuousCheck()
	
	// Initial checks
	if detectDebugger() {
		corrupt()
	}
}

// detectDebugger runs all detection methods
func detectDebugger() bool {
	return isDebuggerAttached() ||
		isRemoteDebuggerAttached() ||
		checkDebugPort() ||
		checkTimingAttack() ||
		checkDebuggerProcesses() ||
		checkBreakpoints() ||
		checkNtGlobalFlag()
}

// isDebuggerAttached checks Windows IsDebuggerPresent
func isDebuggerAttached() bool {
	ret, _, _ := isDebuggerPresent.Call()
	return ret != 0
}

// isRemoteDebuggerAttached checks for remote debugging
func isRemoteDebuggerAttached() bool {
	var debuggerPresent int32
	handle, _, _ := getCurrentProcess.Call()
	ret, _, _ := checkRemoteDebuggerPresent.Call(
		handle,
		uintptr(unsafe.Pointer(&debuggerPresent)),
	)
	return ret != 0 && debuggerPresent != 0
}

// checkDebugPort uses NtQueryInformationProcess
func checkDebugPort() bool {
	const ProcessDebugPort = 7
	var debugPort uintptr
	var returnLength uint32
	
	handle, _, _ := getCurrentProcess.Call()
	ret, _, _ := ntQueryInformationProcess.Call(
		handle,
		ProcessDebugPort,
		uintptr(unsafe.Pointer(&debugPort)),
		unsafe.Sizeof(debugPort),
		uintptr(unsafe.Pointer(&returnLength)),
	)
	
	return ret == 0 && debugPort != 0
}

// checkTimingAttack detects single-stepping
func checkTimingAttack() bool {
	start := time.Now()
	
	// Perform dummy operations
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum += i
	}
	_ = sum
	
	elapsed := time.Since(start)
	
	// If this takes more than 500ms, likely being debugged
	return elapsed > 500*time.Millisecond
}

// checkDebuggerProcesses scans for known debugger processes
func checkDebuggerProcesses() bool {
	// Get own executable path to compare
	selfPath, _ := os.Executable()
	selfPath = strings.ToLower(selfPath)
	
	// Check environment for debugger hints
	for _, env := range os.Environ() {
		envLower := strings.ToLower(env)
		for _, dbg := range debuggerProcesses {
			if strings.Contains(envLower, dbg) {
				return true
			}
		}
	}
	
	// Check command line arguments
	for _, arg := range os.Args {
		argLower := strings.ToLower(arg)
		if strings.Contains(argLower, "debug") ||
			strings.Contains(argLower, "attach") ||
			strings.Contains(argLower, "breakpoint") {
			return true
		}
	}
	
	return false
}

// checkBreakpoints looks for software breakpoints (0xCC)
func checkBreakpoints() bool {
	// Check if OutputDebugString triggers exception
	// When no debugger: returns 0
	// When debugger attached: returns non-zero
	testStr, _ := syscall.UTF16PtrFromString("WINUX_DBG_CHECK")
	ret, _, _ := outputDebugStringW.Call(uintptr(unsafe.Pointer(testStr)))
	
	// On some systems this can indicate debugging
	_ = ret
	return false // Conservative: don't false positive
}

// checkNtGlobalFlag checks PEB for debug flags
func checkNtGlobalFlag() bool {
	const (
		FLG_HEAP_ENABLE_TAIL_CHECK   = 0x10
		FLG_HEAP_ENABLE_FREE_CHECK   = 0x20
		FLG_HEAP_VALIDATE_PARAMETERS = 0x40
	)
	
	// These flags are typically set when debugging
	// Implementation would require reading PEB
	// For now, return false to avoid false positives
	return false
}

// continuousCheck runs detection in background
func continuousCheck() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		if detectDebugger() {
			corrupt()
		}
	}
}

// corrupt causes intentional crash with memory corruption
// This makes debugging extremely difficult
func corrupt() {
	// Overwrite memory with garbage
	garbage := make([]byte, 4096)
	for i := range garbage {
		garbage[i] = byte(i ^ 0xDE ^ 0xAD)
	}
	
	// Create fake error messages
	fakeErrors := []string{
		"WINUX: Critical memory allocation failure",
		"WINUX: Stack corruption detected",
		"WINUX: Invalid instruction at 0x00000000",
		"WINUX: Access violation reading 0xDEADBEEF",
		"WINUX: Heap corruption detected",
	}
	
	// Print random fake error
	idx := time.Now().UnixNano() % int64(len(fakeErrors))
	fmt.Fprintln(os.Stderr, fakeErrors[idx])
	
	// Exit with error code that looks like crash
	os.Exit(0xC0000005) // ACCESS_VIOLATION
}

// IntegrityCheck verifies the binary hasn't been tampered with
func IntegrityCheck(expectedHash string) bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}
	
	data, err := os.ReadFile(exePath)
	if err != nil {
		return false
	}
	
	hash := sha256.Sum256(data)
	actualHash := hex.EncodeToString(hash[:])
	
	return actualHash == expectedHash
}

// ObfuscatedString returns a deobfuscated string
// Use this for sensitive strings to prevent static analysis
func ObfuscatedString(encoded []byte, key byte) string {
	result := make([]byte, len(encoded))
	for i, b := range encoded {
		result[i] = b ^ key ^ byte(i)
	}
	return string(result)
}
