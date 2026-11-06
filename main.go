package main

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	webview2 "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

//go:embed index.html
var page string

const keepBackups = 90

func dataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, "Documents", "Otel Rezervasyon")
}
func dataFile() string  { return filepath.Join(dataDir(), "veriler.json") }
func backupDir() string { return filepath.Join(dataDir(), "yedekler") }

func msg(text string) {
	t, _ := windows.UTF16PtrFromString(text)
	c, _ := windows.UTF16PtrFromString("Otel Rezervasyon")
	windows.MessageBox(0, t, c, windows.MB_OK|windows.MB_ICONWARNING)
}

func loadData() (string, error) {
	b, err := os.ReadFile(dataFile())
	if os.IsNotExist(err) {
		return "", nil
	}
	return string(b), err
}

// Writes to a temp file first then renames, so a power cut can't leave a half-written file.
func saveData(content string) error {
	if err := os.MkdirAll(backupDir(), 0o755); err != nil {
		return err
	}
	tmp := dataFile() + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, dataFile()); err != nil {
		return err
	}
	// one backup per day (overwritten during the day), keep the last N days
	day := filepath.Join(backupDir(), "veriler-"+time.Now().Format("2006-01-02")+".json")
	_ = os.WriteFile(day, []byte(content), 0o644)
	files, _ := filepath.Glob(filepath.Join(backupDir(), "veriler-*.json"))
	sort.Strings(files)
	for len(files) > keepBackups {
		_ = os.Remove(files[0])
		files = files[1:]
	}
	return nil
}

func openFolder() {
	_ = os.MkdirAll(backupDir(), 0o755)
	_ = exec.Command("explorer.exe", dataDir()).Start()
}

func main() {
	// only one copy open at a time, otherwise two windows could overwrite each other
	name, _ := windows.UTF16PtrFromString("OtelRezervasyonTekKopya")
	h, err := windows.CreateMutex(nil, false, name)
	if err == windows.ERROR_ALREADY_EXISTS {
		msg("Program zaten açık. Görev çubuğuna bakın.")
		return
	}
	defer windows.CloseHandle(h)

	_ = os.MkdirAll(dataDir(), 0o755)
	cache, err := os.UserCacheDir()
	if err != nil {
		cache = os.TempDir()
	}

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		AutoFocus: true,
		DataPath:  filepath.Join(cache, "OtelRezervasyon"),
		WindowOptions: webview2.WindowOptions{
			Title:  "Otel Rezervasyon",
			Width:  1400,
			Height: 860,
			Center: true,
		},
	})
	if w == nil {
		msg("Program açılamadı.\n\nBu bilgisayarda Microsoft WebView2 eksik olabilir. " +
			"İnternetten \"WebView2 Runtime indir\" diye aratıp Microsoft'un sitesinden kurun, sonra tekrar deneyin.")
		return
	}
	defer w.Destroy()

	_ = w.Bind("loadData", loadData)
	_ = w.Bind("saveData", saveData)
	_ = w.Bind("openFolder", openFolder)
	_ = w.Bind("dataPath", func() string { return dataDir() })
	w.SetHtml(page)
	w.Run()
}
