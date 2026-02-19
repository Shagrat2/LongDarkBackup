# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Описание проекта

LongDarkBackup — кроссплатформенное Go-приложение для автоматического резервного копирования сохранений игры "The Long Dark". Работает как фоновый процесс (Windows Service / systray), каждые 5 секунд проверяет изменения в файлах сохранений, создаёт бэкапы и предоставляет веб-интерфейс для просмотра и восстановления.

## Сборка

```bash
# Текущая платформа
go build -o LongDarkBackup .

# Windows
GOOS=windows GOARCH=amd64 go build -o LDBackup-amd64.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -o LongDarkBackup .

# macOS
GOOS=darwin GOARCH=amd64 go build -o LongDarkBackup .

# WASM (для веб-версии UI)
GOOS=js GOARCH=wasm go build -o web/app.wasm .
```

Тесты отсутствуют. Makefile и CI/CD нет — используются VS Code tasks (.vscode/tasks.json).

## Архитектура

Весь код в одном пакете `main`. Платформо-зависимая логика разделена через build-теги в именах файлов (`_win`, `_nix`, `_wasm`).

### Основной цикл

`main.go` → `kardianos/service` (Windows Service) + `mainthread` → запуск systray + HTTP-сервер на `localhost:45192` + фоновое сканирование.

### Ключевые модули

| Файл | Назначение |
|------|-----------|
| `main.go` | Точка входа, сервис, systray, определение путей к сохранениям |
| `scan.go` | Фоновое сканирование изменений (каждые 5 сек, ожидание до 30 сек) |
| `fileList.go` | Отслеживание файлов по времени модификации |
| `cache.go` | Парсинг сохранений: LZF-декомпрессия → JSON → метаданные и скриншот |
| `app.go` | HTTP-сервер, маршруты go-app/v10, встроенные статические ресурсы |
| `appList.go` | Веб-UI: список бэкапов, организованных по имени/дате |
| `restore.go` | HTTP-обработчик `/restore/{id}` — восстановление из бэкапа |
| `lockFile_*.go` | Блокировка файлов: Windows (LockFileEx), Unix (flock), WASM (no-op) |
| `systray_*.go` | System tray: полная реализация (не-WASM) / stub (WASM) |
| `statFIle.go` | `//go:embed` — встраивание favicon, CSS, WASM в бинарник |

### Формат сохранений The Long Dark

Файлы `sandbox*` содержат LZF-сжатый JSON с вложенными JSON-строками в `m_Dict["info"]` (статистика) и `m_Dict["screenshot"]` (base64 JPEG).

### Переменные окружения

- `LDB_DATA_DIR` — переопределяет путь к сохранениям игры
- `LDB_BACKUP_DIR` — переопределяет путь к резервным копиям

Стандартные пути: Windows — `%APPDATA%/Hinterland/TheLongDark/Survival`, Linux/macOS — `~/.local/share/Hinterland/TheLongDark/Survival`. Бэкапы: `~/Documents/LongDarkBackup`.

### Ключевые зависимости

- `go-app/v10` — PWA-фреймворк (UI компилируется и в нативный сервер, и в WASM)
- `kardianos/service` — запуск как Windows Service
- `getlantern/systray` — иконка в системном трее
- `zhuyie/golzf` — LZF декомпрессия сохранений
