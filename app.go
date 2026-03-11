package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	// User data path cache
	userDataPath string
	// App path cache
	appPath string
	// Pending update file path
	pendingUpdateFile string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Launcher Methods

// Name returns the launcher name
func (a *App) Name() string {
	return "Keebell"
}

// Version returns the launcher version
func (a *App) Version() string {
	return "2.11.0"
}

// Platform returns the current platform (darwin, win32, linux)
func (a *App) Platform() string {
	switch runtime.GOOS {
	case "darwin":
		return "darwin"
	case "windows":
		return "win32"
	case "linux":
		return "linux"
	default:
		return runtime.GOOS
	}
}

// Arch returns the current architecture
func (a *App) Arch() string {
	return runtime.GOARCH
}

// AutoTypeSupported returns whether auto-type is supported
func (a *App) AutoTypeSupported() bool {
	// Disable auto-type for initial version
	return false
}

// ThirdPartyStoragesSupported returns whether third party storages are supported
func (a *App) ThirdPartyStoragesSupported() bool {
	// Enable basic support
	return true
}

// ClipboardSupported returns whether clipboard is supported
func (a *App) ClipboardSupported() bool {
	return true
}

// DevTools returns whether dev tools are available
func (a *App) DevTools() bool {
	// Enable in development
	return true
}

// OpenDevTools opens the developer tools
func (a *App) OpenDevTools() {
	// DevTools not available in Wails v2
	// wailsruntime.WindowOpenDevTools(a.ctx)
}

// OpenLink opens an external link
func (a *App) OpenLink(href string) {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") ||
		strings.HasPrefix(href, "ftp://") || strings.HasPrefix(href, "sftp://") ||
		strings.HasPrefix(href, "mailto:") {
		wailsruntime.BrowserOpenURL(a.ctx, href)
	}
}

// GetUserDataPath returns the path to user data directory
func (a *App) GetUserDataPath(fileName string) string {
	if a.userDataPath == "" {
		// Get user config directory
		configDir, err := os.UserConfigDir()
		if err != nil {
			// Fallback to current directory
			configDir = "."
		}
		a.userDataPath = filepath.Join(configDir, "keebell")

		// Ensure directory exists
		os.MkdirAll(a.userDataPath, 0755)
	}

	if fileName == "" {
		return a.userDataPath
	}
	return filepath.Join(a.userDataPath, fileName)
}

// GetTempPath returns the path to temp directory
func (a *App) GetTempPath(fileName string) string {
	tempDir := os.TempDir()
	keebellTemp := filepath.Join(tempDir, "keebell")

	// Ensure directory exists
	os.MkdirAll(keebellTemp, 0755)

	if fileName == "" {
		return keebellTemp
	}
	return filepath.Join(keebellTemp, fileName)
}

// GetDocumentsPath returns the path to documents directory
func (a *App) GetDocumentsPath(fileName string) string {
	docDir, err := os.UserHomeDir()
	if err != nil {
		docDir = "."
	}
	docDir = filepath.Join(docDir, "Documents")

	if fileName == "" {
		return docDir
	}
	return filepath.Join(docDir, fileName)
}

// GetAppPath returns the path to application directory
func (a *App) GetAppPath(fileName string) string {
	if a.appPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			a.appPath = "."
		} else {
			a.appPath = filepath.Dir(exePath)
		}
	}

	if fileName == "" {
		return a.appPath
	}
	return filepath.Join(a.appPath, fileName)
}

// GetWorkDirPath returns the path to current working directory
func (a *App) GetWorkDirPath(fileName string) string {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	if fileName == "" {
		return wd
	}
	return filepath.Join(wd, fileName)
}

// JoinPath joins path segments
func (a *App) JoinPath(parts ...string) string {
	return filepath.Join(parts...)
}

// WriteFile writes data to a file
func (a *App) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// ReadFile reads data from a file
func (a *App) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// FileExists checks if a file exists
func (a *App) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DeleteFile deletes a file
func (a *App) DeleteFile(path string) error {
	return os.Remove(path)
}

// StatFile gets file stats
func (a *App) StatFile(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Mkdir creates a directory recursively
func (a *App) Mkdir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// ParsePath parses a file path into components
func (a *App) ParsePath(fileName string) map[string]string {
	return map[string]string{
		"path": fileName,
		"dir":  filepath.Dir(fileName),
		"file": filepath.Base(fileName),
	}
}

// SetClipboardText sets clipboard text
func (a *App) SetClipboardText(text string) error {
	wailsruntime.ClipboardSetText(a.ctx, text)
	return nil
}

// GetClipboardText gets clipboard text
func (a *App) GetClipboardText() string {
	text, err := wailsruntime.ClipboardGetText(a.ctx)
	if err != nil {
		return ""
	}
	return text
}

// ClearClipboardText clears clipboard
func (a *App) ClearClipboardText() {
	// Wails doesn't have direct clipboard clear, set empty
	wailsruntime.ClipboardSetText(a.ctx, "")
}

// ShowMainWindow shows and focuses the main window
func (a *App) ShowMainWindow() {
	wailsruntime.Show(a.ctx)
}

// HideApp hides the application
func (a *App) HideApp() {
	wailsruntime.Hide(a.ctx)
}

// IsAppFocused checks if app window is focused
func (a *App) IsAppFocused() bool {
	// WindowIsFocused not available in Wails v2, assume focused if not minimized
	return !wailsruntime.WindowIsMinimised(a.ctx)
}

// MinimizeApp minimizes to tray (placeholder)
func (a *App) MinimizeApp(restoreText, quitText string) {
	wailsruntime.WindowMinimise(a.ctx)
}

// MinimizeMainWindow minimizes the main window
func (a *App) MinimizeMainWindow() {
	wailsruntime.WindowMinimise(a.ctx)
}

// MaximizeMainWindow maximizes the main window
func (a *App) MaximizeMainWindow() {
	wailsruntime.WindowMaximise(a.ctx)
}

// RestoreMainWindow restores the main window
func (a *App) RestoreMainWindow() {
	wailsruntime.WindowUnmaximise(a.ctx)
}

// MainWindowMaximized checks if main window is maximized
func (a *App) MainWindowMaximized() bool {
	return wailsruntime.WindowIsMaximised(a.ctx)
}

// GetSaveFileName shows save file dialog
func (a *App) GetSaveFileName(defaultPath, title, filterName, filterExt string) string {
	result, err := wailsruntime.SaveFileDialog(a.ctx, wailsruntime.SaveDialogOptions{
		DefaultFilename: defaultPath,
		Title:           title,
		Filters: []wailsruntime.FileFilter{
			{
				DisplayName: filterName,
				Pattern:     filterExt,
			},
		},
	})
	if err != nil {
		return ""
	}
	return result
}

// LoadConfig loads configuration
func (a *App) LoadConfig(name string) (string, error) {
	configPath := filepath.Join(a.GetUserDataPath(""), name+".json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // Return empty for non-existent file
		}
		return "", err
	}
	return string(data), nil
}

// SaveConfig saves configuration
func (a *App) SaveConfig(name, data string) error {
	configPath := filepath.Join(a.GetUserDataPath(""), name+".json")
	return os.WriteFile(configPath, []byte(data), 0644)
}

// Exit requests application exit
func (a *App) Exit() {
	a.pendingUpdateFile = ""
	wailsruntime.Quit(a.ctx)
}

// RequestExit requests exit (for updates)
func (a *App) RequestExit() {
	if a.pendingUpdateFile != "" {
		// TODO: Implement update restart
		fmt.Println("Would restart with update:", a.pendingUpdateFile)
	}
	wailsruntime.Quit(a.ctx)
}

// RequestRestartAndUpdate schedules restart with update
func (a *App) RequestRestartAndUpdate(updateFilePath string) {
	a.pendingUpdateFile = updateFilePath
}

// CancelRestart cancels pending restart
func (a *App) CancelRestart() {
	a.pendingUpdateFile = ""
}

// Spawn spawns a process (placeholder)
func (a *App) Spawn(config map[string]interface{}) map[string]interface{} {
	// TODO: Implement process spawning
	fmt.Println("Spawn called with config:", config)
	return map[string]interface{}{
		"code":   0,
		"stdout": "",
		"stderr": "",
		"err":    nil,
	}
}

// SetGlobalShortcuts sets global shortcuts (placeholder)
func (a *App) SetGlobalShortcuts(appSettings map[string]interface{}) {
	// TODO: Implement global shortcuts
	fmt.Println("SetGlobalShortcuts called with settings:", appSettings)
}

// CheckOpenFiles checks for files to open (placeholder)
func (a *App) CheckOpenFiles() {
	// TODO: Implement file opening queue
}

// OpenFile opens a file
func (a *App) OpenFile(file string) {
	// TODO: Implement file opening logic
	fmt.Println("OpenFile called with:", file)
}

// ResolveProxy resolves proxy for URL (placeholder)
func (a *App) ResolveProxy(url string) map[string]interface{} {
	// TODO: Implement proxy resolution
	return map[string]interface{}{
		"host": "",
		"port": 0,
	}
}

// CanDetectOsSleep checks if OS sleep detection is supported
func (a *App) CanDetectOsSleep() bool {
	return a.Platform() != "linux"
}

// UpdaterEnabled checks if updater is enabled
func (a *App) UpdaterEnabled() bool {
	return a.Platform() != "linux"
}
