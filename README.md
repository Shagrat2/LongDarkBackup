# LongDarkBackup

[Русская версия](README_RU.md)

Automatic backup service for **The Long Dark** saves. Runs as a background process (Windows Service / systray), monitors save files for changes every 5 seconds, creates backups, and provides a web interface to browse and restore them.

## Features

- Automatic backup on save file changes (only changed files are backed up)
- Web interface at `http://localhost:45192` with dark theme in The Long Dark style
- Detailed save info: region, survival time, condition, afflictions, screenshot
- One-click restore from any backup
- Multilingual interface (English / Russian)
- Cross-platform: Windows, Linux, macOS
- System tray icon with quick access

## Installation

Download the latest release from the [releases page](https://github.com/Shagrat2/LongDarkBackup/releases):

| Platform | File | Install |
|----------|------|---------|
| Windows x32 | `.zip` | Unzip and run `LongDarkBackup.exe` |
| Windows x64 | `.zip` | Unzip and run `LongDarkBackup.exe` |
| macOS Intel | `.dmg` | Open DMG, drag to Applications |
| macOS Apple Silicon | `.dmg` | Open DMG, drag to Applications |
| Linux x64 | `.zip` | Unzip and run `./LongDarkBackup` |

> **macOS**: If blocked by Gatekeeper — System Settings → Privacy & Security → Open Anyway.

The program will open the web interface in your browser automatically.

## Build from source

Requires Go 1.23+. Linux build requires Docker (for GTK/systray).

```bash
# Current platform
go build -o LongDarkBackup ./app/

# All platforms release build
./release.sh v1.0.0

# Build + create GitHub release
./release.sh v1.0.0 --publish
```

## Configuration

Environment variables (optional):

| Variable | Description | Default |
|----------|-------------|---------|
| `LDB_DATA_DIR` | Path to game saves | Windows: `%APPDATA%/Hinterland/TheLongDark/Survival`<br>Linux/macOS: `~/.local/share/Hinterland/TheLongDark/Survival` |
| `LDB_BACKUP_DIR` | Path to backups | `~/Documents/LongDarkBackup` |

## How It Works

1. Every 5 seconds the program scans the game save directory for changes
2. If changes are detected, it waits for writes to finish (up to 30 seconds)
3. Only the changed files are copied to a timestamped backup folder
4. Save metadata is extracted from LZF-compressed sandbox files (stats, screenshot, afflictions)
5. The web interface shows all backups with screenshots, game stats, and condition bars

## Dependencies

- [go-app/v10](https://github.com/maxence-charriere/go-app) — PWA framework (SSR mode)
- [kardianos/service](https://github.com/kardianos/service) — Windows Service support
- [getlantern/systray](https://github.com/getlantern/systray) — System tray icon
- [zhuyie/golzf](https://github.com/zhuyie/golzf) — LZF decompression for save files

## License

[Apache License 2.0](LICENSE)

Copyright 2024-2026 Shagrat2. See [NOTICE](NOTICE) for attribution requirements.
