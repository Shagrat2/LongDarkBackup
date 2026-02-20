# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Описание проекта

LongDarkBackup — кроссплатформенное Go-приложение для автоматического резервного копирования сохранений игры "The Long Dark". Работает как фоновый процесс (Windows Service / systray), каждые 5 секунд проверяет изменения в файлах сохранений, создаёт бэкапы и предоставляет веб-интерфейс (EN/RU) для просмотра и восстановления.

## Сборка и релиз

```bash
# Текущая платформа
go build -o LongDarkBackup ./app/

# Релиз всех платформ (Windows zip, macOS dmg, Linux zip)
./release.sh v1.0.0

# Релиз + публикация на GitHub
./release.sh v1.0.0 --publish
```

Linux собирается через Docker (`Dockerfile.linux`) — требуется для GTK/systray.
macOS собирается как .app бандл внутри .dmg с ad-hoc подписью.
Тесты отсутствуют. VS Code tasks в `.vscode/tasks.json`: build, release, release + publish.

## Архитектура

Весь код в пакете `main` в директории `app/`. Платформо-зависимая логика разделена через build-теги в именах файлов (`_win`, `_nix`). Статические ресурсы в `app/web/`.

WASM-сборка не поддерживается — приложение работает только в SSR-режиме через go-app/v10.

### Основной цикл

`app/main.go` → `kardianos/service` (Windows Service) + `mainthread` → запуск systray + HTTP-сервер на `localhost:45192` + фоновое сканирование.

### Ключевые модули

| Файл | Назначение |
|------|-----------|
| `app/main.go` | Точка входа, сервис, systray, определение путей к сохранениям |
| `app/scan.go` | Фоновое сканирование изменений (каждые 5 сек, ожидание до 30 сек), бэкап только изменённых файлов |
| `app/fileList.go` | Отслеживание файлов по времени модификации |
| `app/cache.go` | Парсинг сохранений: LZF-декомпрессия → JSON → метаданные, скриншот, болезни |
| `app/app.go` | HTTP-сервер, маршруты, middleware (lang, noETag), встроенные ресурсы |
| `app/appList.go` | Главная страница: список бэкапов по имени/дате (go-app SSR) |
| `app/appDetail.go` | Детальная страница бэкапа (html/template), восстановление |
| `app/restore.go` | HTTP-обработчик `/restore/{id}` — восстановление из бэкапа |
| `app/i18n.go` | Мультиязычность: переводы EN/RU, определение языка, функция `T()` |
| `app/lockFile_win.go` | Блокировка файлов: Windows (LockFileEx) |
| `app/lockFile_nix.go` | Блокировка файлов: Unix (flock) |
| `app/systray_other.go` | System tray: darwin + windows (CGO, getlantern/systray) |
| `app/systray_linux.go` | System tray: linux (CGO, getlantern/systray через Docker) |
| `app/statFIle.go` | `//go:embed` — встраивание favicon и CSS в бинарник |

### Релизная сборка

| Файл | Назначение |
|------|-----------|
| `release.sh` | Сборка всех платформ + публикация GitHub release |
| `Dockerfile.linux` | Docker-образ для сборки Linux (golang + GTK + ayatana-appindicator) |
| `pkg/macos/Info.plist` | .app бандл: CFBundle, LSUIElement=true (без Dock-иконки) |
| `pkg/macos/icon.icns` | Иконка macOS .app |

### i18n (мультиязычность)

- `app/i18n.go` — все строки перевода, функция `T(key)`, определение языка из cookie/Accept-Language
- Язык хранится в глобальной переменной `currentLang` (безопасно для однопользовательского localhost)
- Переключатель RU/EN на каждой странице, язык сохраняется в cookie (1 год) + localStorage
- `langMiddleware` в `app/app.go` устанавливает язык до рендеринга

### go-app SSR без WASM

go-app/v10 используется только в SSR-режиме. WASM отключён:
- `/web/app.wasm` возвращает 404
- Загрузочный экран go-app скрыт через CSS (`#app-wasm-loader{display:none}`)
- `noETag` middleware убирает `If-None-Match` для корректного переключения языка

### Формат сохранений The Long Dark

Файлы `sandbox*` содержат LZF-сжатый JSON с вложенными JSON-строками в `m_Dict["info"]` (статистика) и `m_Dict["screenshot"]` (base64 JPEG). Глобальные данные (болезни, голод) в `m_Dict["global"]`.

### Переменные окружения

- `LDB_DATA_DIR` — переопределяет путь к сохранениям игры
- `LDB_BACKUP_DIR` — переопределяет путь к резервным копиям

Стандартные пути: Windows — `%APPDATA%/Hinterland/TheLongDark/Survival`, Linux/macOS — `~/.local/share/Hinterland/TheLongDark/Survival`. Бэкапы: `~/Documents/LongDarkBackup`.

### Ключевые зависимости

- `go-app/v10` — PWA-фреймворк (SSR-режим, без WASM)
- `kardianos/service` — запуск как Windows Service
- `getlantern/systray` — иконка в системном трее (CGO на Linux/macOS)
- `zhuyie/golzf` — LZF декомпрессия сохранений
