//go:build windows
// +build windows

// Windows Update Notification Dialog
// Modern MessageBox-style popup for update notifications

package updater

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	procMessageBoxW         = user32.NewProc("MessageBoxW")
	procGetConsoleWindow    = kernel32.NewProc("GetConsoleWindow")
	procCreateWindowExW     = user32.NewProc("CreateWindowExW")
	procDefWindowProcW      = user32.NewProc("DefWindowProcW")
	procRegisterClassExW    = user32.NewProc("RegisterClassExW")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procPostQuitMessage     = user32.NewProc("PostQuitMessage")
	procShowWindow          = user32.NewProc("ShowWindow")
	procUpdateWindow        = user32.NewProc("UpdateWindow")
	procDestroyWindow       = user32.NewProc("DestroyWindow")
	procSetWindowTextW      = user32.NewProc("SetWindowTextW")
	procGetModuleHandleW    = kernel32.NewProc("GetModuleHandleW")
	procLoadIconW           = user32.NewProc("LoadIconW")
	procLoadCursorW         = user32.NewProc("LoadCursorW")
	procBeginPaint          = user32.NewProc("BeginPaint")
	procEndPaint            = user32.NewProc("EndPaint")
	procGetClientRect       = user32.NewProc("GetClientRect")
	procDrawTextW           = user32.NewProc("DrawTextW")
	procCreateFontW         = user32.NewProc("CreateFontW")
	procSelectObject        = user32.NewProc("SelectObject")
	procDeleteObject        = user32.NewProc("DeleteObject")
	procSetBkMode           = user32.NewProc("SetBkMode")
	procSetTextColor        = user32.NewProc("SetTextColor")
	procGetStockObject      = user32.NewProc("GetStockObject")
	procSendMessageW        = user32.NewProc("SendMessageW")
	procSetFocus            = user32.NewProc("SetFocus")

	gdi32            = syscall.NewLazyDLL("gdi32.dll")
	procCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
)

const (
	// MessageBox Constants
	MB_OK              = 0x00000000
	MB_OKCANCEL        = 0x00000001
	MB_YESNO           = 0x00000004
	MB_YESNOCANCEL     = 0x00000003
	MB_ICONINFORMATION = 0x00000040
	MB_ICONWARNING     = 0x00000030
	MB_ICONQUESTION    = 0x00000020
	MB_DEFBUTTON1      = 0x00000000
	MB_DEFBUTTON2      = 0x00000100
	MB_SETFOREGROUND   = 0x00010000
	MB_TOPMOST         = 0x00040000

	// Button IDs
	IDOK     = 1
	IDCANCEL = 2
	IDYES    = 6
	IDNO     = 7
	
	// Custom button IDs
	ID_BTN_UPDATE    = 1001
	ID_BTN_LATER     = 1002
	ID_BTN_MORE      = 1003

	// Window styles
	WS_OVERLAPPED       = 0x00000000
	WS_CAPTION          = 0x00C00000
	WS_SYSMENU          = 0x00080000
	WS_MINIMIZEBOX      = 0x00020000
	WS_VISIBLE          = 0x10000000
	WS_CHILD            = 0x40000000
	WS_TABSTOP          = 0x00010000
	WS_EX_DLGMODALFRAME = 0x00000001
	WS_EX_TOPMOST       = 0x00000008

	// Button styles
	BS_PUSHBUTTON  = 0x00000000
	BS_DEFPUSHBUTTON = 0x00000001

	// Window messages
	WM_CREATE    = 0x0001
	WM_DESTROY   = 0x0002
	WM_CLOSE     = 0x0010
	WM_PAINT     = 0x000F
	WM_COMMAND   = 0x0111
	WM_SETFONT   = 0x0030

	// Colors
	COLOR_WINDOW = 5
	COLOR_BTNFACE = 15

	// Draw text
	DT_LEFT       = 0x00000000
	DT_WORDBREAK  = 0x00000010
	DT_NOPREFIX   = 0x00000800

	// Background mode
	TRANSPARENT = 1

	// Stock objects
	DEFAULT_GUI_FONT = 17

	// Icons
	IDI_INFORMATION = 32516
	IDC_ARROW       = 32512

	// ShowWindow
	SW_SHOW = 5
)

// RECT structure
type RECT struct {
	Left, Top, Right, Bottom int32
}

// PAINTSTRUCT structure
type PAINTSTRUCT struct {
	HDC         uintptr
	FErase      int32
	RcPaint     RECT
	FRestore    int32
	FIncUpdate  int32
	RgbReserved [32]byte
}

// MSG structure
type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

// WNDCLASSEX structure
type WNDCLASSEX struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}

// Dialog state
var (
	dialogResult    int
	updateInfo      *UpdateInfo
	hFont           uintptr
	btnUpdate       uintptr
	btnLater        uintptr
	btnMore         uintptr
)

// ShowUpdateNotification displays a Windows dialog with update information.
// Returns: 0 = Later, 1 = Update Now, 2 = More Info
func ShowUpdateNotification(info *UpdateInfo) int {
	if !info.Available {
		return 0
	}

	updateInfo = info
	dialogResult = 0

	// Build the message with changelog summary
	var msgBuilder strings.Builder
	msgBuilder.WriteString(fmt.Sprintf("🚀 WINUX Güncellemesi Mevcut!\n\n"))
	msgBuilder.WriteString(fmt.Sprintf("Mevcut Sürüm: v%s\n", info.CurrentVersion))
	msgBuilder.WriteString(fmt.Sprintf("Yeni Sürüm: v%s\n\n", info.LatestVersion))
	msgBuilder.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	msgBuilder.WriteString("📋 Neler Yeni:\n\n")
	
	for i, line := range info.Summary {
		if i >= 5 {
			break
		}
		msgBuilder.WriteString(line + "\n")
	}
	
	msgBuilder.WriteString("\n━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	msgBuilder.WriteString("\n🔄 Şimdi güncellemek ister misiniz?")

	message := msgBuilder.String()
	title := "WINUX Update"

	// Use simple MessageBox for reliable display
	result := messageBoxCustom(title, message, info.ReleaseURL)
	
	return result
}

// messageBoxCustom shows a simple Yes/No/More dialog
func messageBoxCustom(title, message, releaseURL string) int {
	// For simplicity, use standard MessageBox with Yes/No
	// Yes = Update, No = Later
	// We'll open release URL separately if they want more info
	
	fullMessage := message + "\n\n[Evet] = Güncelle | [Hayır] = Sonra"
	
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	msgPtr, _ := syscall.UTF16PtrFromString(fullMessage)
	
	ret, _, _ := procMessageBoxW.Call(
		0,
		uintptr(unsafe.Pointer(msgPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(MB_YESNOCANCEL|MB_ICONINFORMATION|MB_SETFOREGROUND|MB_TOPMOST|MB_DEFBUTTON1),
	)
	
	switch ret {
	case IDYES:
		return 1 // Update
	case IDNO:
		return 0 // Later
	case IDCANCEL:
		// Cancel = More Info, open browser
		if releaseURL != "" {
			OpenURL(releaseURL)
		}
		return 2 // More Info
	default:
		return 0
	}
}

// ShowUpdateAvailableToast shows a simple toast notification (fallback)
func ShowUpdateAvailableToast(info *UpdateInfo) {
	if !info.Available {
		return
	}
	
	title := "WINUX Update Available"
	message := fmt.Sprintf("Version %s is available! Run 'update.exe --apply' to update.", info.LatestVersion)
	
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	msgPtr, _ := syscall.UTF16PtrFromString(message)
	
	procMessageBoxW.Call(
		0,
		uintptr(unsafe.Pointer(msgPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(MB_OK|MB_ICONINFORMATION|MB_SETFOREGROUND|MB_TOPMOST),
	)
}

// utf16Ptr converts a string to UTF-16 pointer
func utf16Ptr(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}
