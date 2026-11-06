# Otel Rezervasyon

A tiny offline room-booking manager for a small boutique hotel, built to run on a single Windows PC and to be usable by someone who isn't comfortable with technology.

One `.exe`, no installer, no internet, no database server. The UI is in Turkish.

## Features

- **Calendar grid**: rooms as rows, days of the month as columns. Click an empty cell to book, click a booking to edit or delete it.
- **Double-booking protection**: overlapping stays in the same room are rejected.
- **Payment status at a glance**: green = paid, orange = deposit taken, red = unpaid. The form shows remaining balance and number of nights.
- **Today panel**: arrivals, departures, and guests currently in house.
- **Search** by guest name, phone, or note.
- **Editable room list**: rooms have names rather than numbers.
- **Safe storage**:
  - Data is written to `Documents\Otel Rezervasyon\veriler.json` on every save.
  - Writes are atomic (temp file + rename), so a power cut can't corrupt the file.
  - One backup per day goes to `yedekler\`, and the last 90 days are kept.
  - Single-instance lock, so two windows can't overwrite each other.
  - If the data file can't be read, the app refuses to overwrite it.

## How it works

| Layer | Choice | Why |
|---|---|---|
| Window | [go-webview2](https://github.com/jchv/go-webview2) | Uses the WebView2 runtime that ships with Windows 10/11. Pure Go, no CGO, so it cross-compiles from Linux/macOS. |
| UI | Single `index.html` (vanilla JS/CSS), embedded with `go:embed` | No build step, no framework. |
| Storage | JSON file, accessed through functions bound from Go into JS (`loadData`, `saveData`, `openFolder`) | Human-readable, trivial to back up, no DB to install. |

The result is a ~3 MB executable.

## Build

You need Go 1.24 or newer. You can build from any OS:

```bash
go mod tidy
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-H windowsgui -s -w" -o OtelRezervasyon.exe .
```

On Windows (PowerShell):

```powershell
go mod tidy
$env:GOOS="windows"; $env:CGO_ENABLED="0"
go build -trimpath -ldflags "-H windowsgui -s -w" -o OtelRezervasyon.exe .
```

GitHub Actions builds the exe on every push to `main` and uploads it as a workflow artifact. Pushing a `v*` tag also attaches the exe to a GitHub Release:

```bash
git tag v1.0.0 && git push origin v1.0.0
```

## Running

1. Download `OtelRezervasyon.exe` from Releases and put it on the desktop.
2. Double-click it. The exe is unsigned, so Windows SmartScreen may warn on the first run. Click **More info → Run anyway**.
3. Click **Odalar** to set your room names.

If WebView2 is missing (rare, only on very old, un-updated Windows 10 machines), the app says so. Install the WebView2 Runtime from Microsoft and run it again.

## Data & backups

```
Documents\Otel Rezervasyon\
├── veriler.json          ← live data
└── yedekler\
    ├── veriler-2026-11-05.json
    └── ...               ← last 90 days
```

To restore an old backup, click **Yedekten Yükle** in the app and pick a file from `yedekler\`. To move to a new PC, copy the whole `Otel Rezervasyon` folder.

## Roadmap

- Public booking website with online payment, writing into the same data
- Price per room/season
- Export to Excel
