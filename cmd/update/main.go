// WINUX Update - Self-update utility for WINUX
// Downloads and applies updates from GitHub Releases
//
// Usage: update.exe [--check | --apply | --force]
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// GitHub API endpoint
	GitHubAPI      = "https://api.github.com/repos/CRTYPUBG/winux/releases/latest"
	CurrentVersion = "0.1.0"
	BinaryName     = "winux.exe"
	UpdaterName    = "update.exe"
)

// GitHubRelease represents the GitHub API response
type GitHubRelease struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Body    string  `json:"body"`
	Assets  []Asset `json:"assets"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "--check", "-c":
		checkForUpdates()
	case "--apply", "-a":
		applyUpdate(false)
	case "--force", "-f":
		applyUpdate(true)
	case "--help", "-h":
		printUsage()
	case "--version", "-v":
		fmt.Printf("WINUX Updater v%s\n", CurrentVersion)
	default:
		fmt.Fprintf(os.Stderr, "Unknown option: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`WINUX Update - Self-update utility

Usage: update.exe [option]

Options:
  --check, -c     Check for available updates
  --apply, -a     Download and apply update if available
  --force, -f     Force reinstall even if up-to-date
  --version, -v   Show version
  --help, -h      Show this help

Examples:
  update.exe --check
  update.exe --apply`)
}

func checkForUpdates() {
	fmt.Println("Checking for updates...")
	
	release, err := getLatestRelease()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	
	if compareVersions(latestVersion, CurrentVersion) > 0 {
		fmt.Printf("\n✅ Update available!\n")
		fmt.Printf("   Current: v%s\n", CurrentVersion)
		fmt.Printf("   Latest:  %s\n", release.TagName)
		fmt.Printf("\nRun 'update.exe --apply' to update.\n")
	} else {
		fmt.Printf("\n✓ WINUX is up-to-date (v%s)\n", CurrentVersion)
	}
}

func applyUpdate(force bool) {
	fmt.Println("WINUX Update")
	fmt.Println("============")

	// Get latest release info
	release, err := getLatestRelease()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching release info: %v\n", err)
		os.Exit(1)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")

	if !force && compareVersions(latestVersion, CurrentVersion) <= 0 {
		fmt.Printf("Already up-to-date (v%s)\n", CurrentVersion)
		return
	}

	fmt.Printf("Updating to %s...\n\n", release.TagName)

	// Find binary and checksum assets
	var binaryURL, checksumURL string
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".exe") && !strings.Contains(asset.Name, "update") {
			binaryURL = asset.BrowserDownloadURL
		}
		if strings.HasSuffix(asset.Name, ".sha256") {
			checksumURL = asset.BrowserDownloadURL
		}
	}

	if binaryURL == "" {
		fmt.Fprintln(os.Stderr, "Error: Binary not found in release assets")
		os.Exit(1)
	}

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "winux-update-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	newBinaryPath := filepath.Join(tempDir, BinaryName)

	// Download new binary
	fmt.Printf("[1/4] Downloading %s...\n", release.TagName)
	if err := downloadFile(binaryURL, newBinaryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("      ✓ Download complete")

	// Verify checksum if available
	if checksumURL != "" {
		fmt.Println("[2/4] Verifying SHA256 checksum...")
		checksumPath := filepath.Join(tempDir, "checksum.sha256")
		if err := downloadFile(checksumURL, checksumPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not download checksum: %v\n", err)
		} else {
			if err := verifyChecksum(newBinaryPath, checksumPath); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Checksum verification failed: %v\n", err)
				fmt.Fprintln(os.Stderr, "Update aborted for security reasons.")
				os.Exit(1)
			}
			fmt.Println("      ✓ Checksum verified")
		}
	} else {
		fmt.Println("[2/4] Skipping checksum (not available)")
	}

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding current executable: %v\n", err)
		os.Exit(1)
	}
	currentDir := filepath.Dir(currentExe)
	targetPath := filepath.Join(currentDir, BinaryName)

	// Backup old binary
	fmt.Println("[3/4] Backing up current version...")
	backupPath := targetPath + ".backup"
	if _, err := os.Stat(targetPath); err == nil {
		if err := os.Rename(targetPath, backupPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error backing up: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("      ✓ Backup created")
	} else {
		fmt.Println("      - No existing binary to backup")
	}

	// Install new binary
	fmt.Println("[4/4] Installing new version...")
	if err := copyFile(newBinaryPath, targetPath); err != nil {
		// Rollback on failure
		if _, backupErr := os.Stat(backupPath); backupErr == nil {
			os.Rename(backupPath, targetPath)
		}
		fmt.Fprintf(os.Stderr, "Error installing: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("      ✓ Installation complete")

	// Remove backup
	os.Remove(backupPath)

	fmt.Printf("\n✅ Successfully updated to %s!\n", release.TagName)
	fmt.Println("\nRestart WINUX to use the new version.")
}

func getLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	
	req, err := http.NewRequest("GET", GitHubAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "WINUX-Updater/"+CurrentVersion)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &release, nil
}

func downloadFile(url, dest string) error {
	client := &http.Client{Timeout: 5 * time.Minute}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "WINUX-Updater/"+CurrentVersion)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Show progress
	counter := &progressWriter{Total: resp.ContentLength}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	fmt.Println() // New line after progress

	return err
}

type progressWriter struct {
	Total      int64
	Downloaded int64
	LastPrint  time.Time
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Downloaded += int64(n)
	
	// Print progress every 100ms
	if time.Since(pw.LastPrint) > 100*time.Millisecond {
		if pw.Total > 0 {
			pct := float64(pw.Downloaded) / float64(pw.Total) * 100
			fmt.Printf("\r      Progress: %.1f%% (%s / %s)", 
				pct, formatBytes(pw.Downloaded), formatBytes(pw.Total))
		} else {
			fmt.Printf("\r      Downloaded: %s", formatBytes(pw.Downloaded))
		}
		pw.LastPrint = time.Now()
	}
	
	return n, nil
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func verifyChecksum(filePath, checksumPath string) error {
	// Read expected checksum
	checksumData, err := os.ReadFile(checksumPath)
	if err != nil {
		return err
	}
	
	// Parse checksum file (format: "hash  filename")
	parts := strings.Fields(string(checksumData))
	if len(parts) < 1 {
		return fmt.Errorf("invalid checksum file format")
	}
	expectedHash := strings.ToLower(parts[0])

	// Calculate actual hash
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return err
	}
	actualHash := hex.EncodeToString(hasher.Sum(nil))

	if actualHash != expectedHash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}

func compareVersions(v1, v2 string) int {
	// Simple version comparison (x.y.z)
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := 0; i < 3; i++ {
		var n1, n2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
		}
		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}
	return 0
}

// LaunchAndExit launches the new binary and exits the updater
func LaunchAndExit(binaryPath string, args ...string) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	os.Exit(0)
}
