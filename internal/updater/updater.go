// Package updater provides background update checking with Windows notifications.
package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// Version variables - set at build time via ldflags
var (
	CurrentVersion = "dev" // -X github.com/CRTYPUBG/winux/internal/updater.CurrentVersion=x.y.z
	CheckDelay     = 3 * time.Second
)

const (
	GitHubAPI = "https://api.github.com/repos/CRTYPUBG/winux/releases/latest"
)

// GitHubRelease represents the GitHub API response
type GitHubRelease struct {
	TagName    string  `json:"tag_name"`
	Name       string  `json:"name"`
	Body       string  `json:"body"`
	HTMLURL    string  `json:"html_url"`
	Assets     []Asset `json:"assets"`
	PublishedAt string `json:"published_at"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// UpdateInfo contains update information
type UpdateInfo struct {
	Available      bool
	CurrentVersion string
	LatestVersion  string
	ReleaseNotes   string        // Full release notes
	Summary        []string      // First 5 lines of changelog
	ReleaseURL     string        // GitHub release page URL
	DownloadURL    string        // Direct download URL
}

// CheckForUpdatesAsync checks for updates in the background with a delay.
// Returns a channel that will receive UpdateInfo when check is complete.
func CheckForUpdatesAsync() <-chan *UpdateInfo {
	ch := make(chan *UpdateInfo, 1)
	
	go func() {
		// Wait before checking (so main app starts first)
		time.Sleep(CheckDelay)
		
		info := CheckForUpdates()
		ch <- info
		close(ch)
	}()
	
	return ch
}

// CheckForUpdates performs a synchronous update check.
func CheckForUpdates() *UpdateInfo {
	info := &UpdateInfo{
		CurrentVersion: CurrentVersion,
	}
	
	release, err := getLatestRelease()
	if err != nil {
		return info
	}
	
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	info.LatestVersion = latestVersion
	info.ReleaseURL = release.HTMLURL
	info.ReleaseNotes = release.Body
	
	// Parse first 5 meaningful lines from changelog
	info.Summary = parseChangelogSummary(release.Body, 5)
	
	// Find download URL
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".exe") && !strings.Contains(asset.Name, "update") {
			info.DownloadURL = asset.BrowserDownloadURL
			break
		}
	}
	
	// Compare versions
	if compareVersions(latestVersion, CurrentVersion) > 0 {
		info.Available = true
	}
	
	return info
}

// parseChangelogSummary extracts the first N meaningful lines from release notes
func parseChangelogSummary(body string, maxLines int) []string {
	lines := strings.Split(body, "\n")
	var summary []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and headers
		if line == "" {
			continue
		}
		
		// Skip markdown headers (##, ###, etc.)
		if strings.HasPrefix(line, "#") && strings.Contains(line, " ") {
			// Keep section headers but clean them
			cleaned := strings.TrimLeft(line, "# ")
			if cleaned != "" && len(summary) < maxLines {
				summary = append(summary, "📋 "+cleaned)
			}
			continue
		}
		
		// Clean up bullet points and add to summary
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
			line = strings.TrimLeft(line, "-* ")
			if line != "" && len(summary) < maxLines {
				summary = append(summary, "• "+line)
			}
		}
		
		if len(summary) >= maxLines {
			break
		}
	}
	
	return summary
}

func getLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	
	req, err := http.NewRequest("GET", GitHubAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "WINUX/"+CurrentVersion)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}
	
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	
	return &release, nil
}

func compareVersions(v1, v2 string) int {
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

// OpenURL opens a URL in the default browser
func OpenURL(url string) error {
	return exec.Command("cmd", "/c", "start", url).Start()
}
