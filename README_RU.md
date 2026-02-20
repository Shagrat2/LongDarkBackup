# LongDarkBackup

[English version](README.md)

Сервис автоматического резервного копирования сохранений **The Long Dark**. Работает как фоновый процесс (Windows Service / systray), каждые 5 секунд проверяет изменения в файлах сохранений, создаёт бэкапы и предоставляет веб-интерфейс для просмотра и восстановления.

## Возможности

- Автоматическое резервное копирование при изменении сохранений (только изменённые файлы)
- Веб-интерфейс на `http://localhost:45192` с тёмной темой в стилистике The Long Dark
- Детальная информация: регион, время выживания, состояние, болезни, скриншот
- Восстановление из любого бэкапа одной кнопкой
- Мультиязычный интерфейс (русский / английский)
- Кроссплатформенность: Windows, Linux, macOS
- Работает как Windows Service или автономно с иконкой в системном трее

## Установка

Скачайте последний релиз со [страницы релизов](https://github.com/Shagrat2/LongDarkBackup/releases), распакуйте и запустите `LongDarkBackup` (или `LDBackup-amd64.exe` на Windows).

Программа автоматически откроет веб-интерфейс в браузере.

## Сборка

Требуется Go 1.21+.

```bash
# Текущая платформа
go build -o LongDarkBackup ./app/

# Windows
GOOS=windows GOARCH=amd64 go build -o LDBackup-amd64.exe ./app/

# Linux
GOOS=linux GOARCH=amd64 go build -o LongDarkBackup ./app/

# macOS
GOOS=darwin GOARCH=amd64 go build -o LongDarkBackup ./app/
```

## Настройка

Переменные окружения (опционально):

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `LDB_DATA_DIR` | Путь к сохранениям игры | Windows: `%APPDATA%/Hinterland/TheLongDark/Survival`<br>Linux/macOS: `~/.local/share/Hinterland/TheLongDark/Survival` |
| `LDB_BACKUP_DIR` | Путь к резервным копиям | `~/Documents/LongDarkBackup` |

## Как это работает

1. Каждые 5 секунд программа сканирует директорию сохранений на наличие изменений
2. При обнаружении изменений ожидает завершения записи (до 30 секунд)
3. Только изменённые файлы копируются в папку бэкапа с временной меткой
4. Метаданные извлекаются из LZF-сжатых sandbox-файлов (статистика, скриншот, болезни)
5. Веб-интерфейс отображает все бэкапы со скриншотами, статистикой и полосками состояния

## Зависимости

- [go-app/v10](https://github.com/maxence-charriere/go-app) — PWA-фреймворк (SSR-режим)
- [kardianos/service](https://github.com/kardianos/service) — поддержка Windows Service
- [getlantern/systray](https://github.com/nicoh-dev/systray) — иконка в системном трее
- [zhuyie/golzf](https://github.com/zhuyie/golzf) — LZF-декомпрессия сохранений

## Лицензия

[Apache License 2.0](LICENSE)

Copyright 2024-2026 Shagrat2. Требования к атрибуции см. в [NOTICE](NOTICE).
