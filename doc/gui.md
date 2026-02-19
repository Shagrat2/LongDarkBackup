# Веб-интерфейс LongDarkBackup

## Визуальный стиль

Тёмная тема в стилистике The Long Dark.

### Палитра

| Роль | Hex |
|------|-----|
| Фон основной | `#0a0e17` |
| Фон карточки | `rgba(20, 28, 40, 0.85)` |
| Текст основной | `#c8d6e5` |
| Текст приглушённый | `#5e6e82` |
| Текст неактивный | `#3a4a5e` |
| Акцент тёплый (янтарный) | `#d4a857` |
| Акцент холодный (ледяной) | `#4fa4b8` |
| Опасность | `#a04040` |
| Хорошее состояние | `#4a8c5c` |
| Предупреждение | `#d4a857` |
| Граница | `rgba(100, 120, 140, 0.3)` |

### Типографика

- **Заголовки**: `Bebas Neue` — condensed sans-serif, uppercase, letter-spacing 2-3px
- **Текст**: `Inter` (wght 300/400/500), fallback `-apple-system, sans-serif`
- **Моноширинный** (время, числа): `ui-monospace, 'SF Mono', monospace`
- Google Fonts: `https://fonts.googleapis.com/css2?family=Bebas+Neue&family=Inter:wght@300;400;500&display=swap`

### Эффекты

- `backdrop-filter: blur(10px)` на карточках
- Бордеры 1px, без теней, без скруглений (резкие формы TLD)
- Transition 0.3s на hover
- Hover карточки: янтарная граница + `scale(1.02)`

---

## Страницы и маршруты

| Маршрут | Файл | Описание |
|---------|------|----------|
| `/` | `app/appList.go` | Главная — список бэкапов (go-app SSR) |
| `/save/{id}` | `app/appDetail.go` | Детальная страница бэкапа (HTML template) |
| `/restore/{id}` | `app/restore.go` | POST → JSON `{"ok":true}`, GET → redirect |
| `/img/{id}` | `app/appDetail.go` | JPEG-скриншот бэкапа |
| `/web/style.css` | `app/statFIle.go` | Стили (embed) |
| `/web/favicon.png` | `app/statFIle.go` | Иконка (embed) |
| `/web/app.wasm` | `app/statFIle.go` | WASM-модуль (embed) |

### Кеширование

| Ресурс | Cache-Control |
|--------|---------------|
| HTML-страницы (`/`, `/save/`) | `no-cache, no-store, must-revalidate` |
| CSS, WASM | `no-cache, no-store, must-revalidate` |
| Скриншоты (`/img/`) | `public, max-age=31536000, immutable` |
| Favicon | `max-age=86400` |

---

## Главная страница (`/`)

Рендерится через go-app/v10 (SSR). Файл: `app/appList.go`.

### Структура данных

```go
type TItem struct {
    Info cacheItem  // метаданные из sandbox-файла
    Path string     // имя папки бэкапа (YYYY-MM-DD-HH-MM-SS)
}

type TDateGroup struct {
    Date      string   // ISO дата (2026-02-19)
    DateLabel string   // русская дата (19 февраля 2026)
    Items     []*TItem
}

type TSaveGroup struct {
    Name  string       // DisplayName из сохранения
    Items []*TItem     // все бэкапы этого сохранения
    Days  []TDateGroup // группировка по датам
}
```

### Группировка и сортировка

1. Все бэкапы сортируются по `Path` в обратном порядке (новые первые)
2. Группируются по `DisplayName` → `TSaveGroup`
3. Внутри каждой группы — подгруппы по дате → `TDateGroup`

### Макет

```
┌──────────────────────────────────────────────────────────────┐
│  LONGDARK BACKUP                                             │
│                                                              │
│  ▶ STALKER  (Voyageur · LakeRegion)              57 бэкапов │
│  ─────────────────────────────────────────────────────────── │
│                                                              │
│  ▼ СТРАДАНИЯ  (Misery · LakeRegion)              32 бэкапа  │
│  ─────────────────────────────────────────────────────────── │
│                                                              │
│    ▶ 19 февраля 2026                             12 бэкапов │
│    ─────────────────────────────────────────────────────     │
│                                                              │
│    ▼ 18 февраля 2026                             20 бэкапов │
│    ─────────────────────────────────────────────────────     │
│    ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐     │
│    │ скриншот │ │ скриншот │ │ скриншот │ │ скриншот │     │
│    │ 19:11    │ │ 19:04    │ │ 18:55    │ │ 18:54    │     │
│    │ Region   │ │ Region   │ │ Region   │ │ Region   │     │
│    │ День 47  │ │ День 47  │ │ День 47  │ │ День 46  │     │
│    │ ██░ 78%  │ │ ██░ 75%  │ │ ██░ 80%  │ │ ██░ 82%  │     │
│    └──────────┘ └──────────┘ └──────────┘ └──────────┘     │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### Компоненты

#### Header

CSS-класс: `.header`

- `h1` с текстом "LONGDARK BACKUP" (Bebas Neue, uppercase, letter-spacing 3px)
- Нижняя граница `1px solid rgba(100, 120, 140, 0.2)`

#### Секция сохранения

HTML: `<details class="save-section">` / `<summary class="save-header">`

- Имя сохранения (`.save-name`) — Bebas Neue, 1.3rem, стрелка ▶ через `::before`
- Подзаголовок (`.save-subtitle`) — режим + регион последнего бэкапа
- Счётчик (`.save-meta`) — "N бэкапов/бэкапа/бэкап"
- Стрелка поворачивается на 90° при `[open]`
- Если одно сохранение — секция раскрыта (`open` атрибут)

#### Секция даты

HTML: `<details class="date-section">` / `<summary class="date-separator">`

- Дата (`.date-text`) — "19 февраля 2026" (Bebas Neue, 1.1rem)
- Счётчик (`.date-count`)
- Свёрнуты по умолчанию
- Стрелка ▶ → ▼ при раскрытии

#### Сетка карточек

CSS-класс: `.card-grid`

- `display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr))`
- Gap: 1rem
- Адаптивность: 768px → minmax(160px), 480px → `1fr 1fr`

#### Карточка бэкапа

HTML: `<a class="card" href="/save/{id}">`

| Элемент | CSS-класс | Содержимое |
|---------|-----------|------------|
| Скриншот | `.card-img` | `<img src="/img/{id}" loading="lazy">`, 130px высота |
| Время бэкапа | `.card-time` | `HH:MM` (ледяной голубой `#4fa4b8`) |
| Регион | `.card-region` | Название региона (обрезка `text-overflow: ellipsis`) |
| Дни выживания | `.card-survived` | "N days HH:MM:SS" |
| Состояние | `.card-condition` | Condition bar + процент |

#### Condition Bar

CSS-классы: `.condition-bar`, `.condition-fill`

- Размер: 100px x 6px (в карточке — flex: 1)
- Фон: `rgba(255, 255, 255, 0.1)`
- Цвет заливки по значению condition:
  - `> 60%` → `.cond-good` → `#4a8c5c` (зелёный)
  - `30-60%` → `.cond-warn` → `#d4a857` (янтарный)
  - `< 30%` → `.cond-danger` → `#a04040` (красный)

---

## Детальная страница (`/save/{id}`)

Рендерится через `html/template`. Файл: `app/appDetail.go`.

### Макет

```
┌──────────────────────────────────────────────────────────────┐
│  ← НАЗАД                                                     │
│                                                              │
│  ┌─────────────────────────┐  MYSTERY LAKE                   │
│  │                         │  Voyageur · SANDBOX              │
│  │                         │                                  │
│  │    скриншот (480x270)   │  Выживание    31 days 23:03:26  │
│  │                         │  Состояние    51.81%  ████░░░   │
│  │                         │  Мир          0.18%              │
│  │                         │  Персонаж     Female             │
│  └─────────────────────────┘                                  │
│                                                              │
│  Бэкап: 19 февраля 2026, 19:04:29                            │
│  Сохранение: СТРАДАНИЯ · 2026-02-19 18:49:22                 │
│                                                              │
│  [ ↻ ВОССТАНОВИТЬ ]                                          │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### Данные шаблона

```go
type DetailData struct {
    Info          cacheItem  // все поля из sandbox
    BackupDateFmt string     // "19 февраля 2026"
    BackupTimeFmt string     // "19:04:29"
    CondClass     string     // "cond-good" / "cond-warn" / "cond-danger"
    CondWidth     string     // "52%"
    ID            string     // имя папки бэкапа
}
```

### Компоненты

| Элемент | CSS-класс | Описание |
|---------|-----------|----------|
| Ссылка назад | `.back-link` | "← НАЗАД", ведёт на `/` |
| Контейнер | `.detail-content` | Flex, gap 2rem |
| Скриншот | `.detail-img` | 480x270, `/img/{id}` |
| Заголовок | `.detail-title` | Регион (Bebas Neue, 1.6rem) |
| Режим | `.detail-mode` | "XPMode · GameMode" |
| Строка статов | `.detail-stat` | Flex space-between, граница снизу |
| Метаданные | `.detail-meta` | Дата бэкапа, время сохранения |
| Кнопка | `.btn-restore` | Янтарная, вызывает `showModal()` |

### Адаптивность

- `< 768px`: `.detail-content` → `flex-direction: column`, скриншот на всю ширину

---

## Модальный диалог восстановления

Показывается на детальной странице при нажатии "ВОССТАНОВИТЬ".

### Макет

```
┌──────────────────────────────────────┐
│                                      │
│  ВОССТАНОВИТЬ СОХРАНЕНИЕ?            │
│                                      │
│  Mystery Lake · 31 days 23:03:26     │
│  Бэкап от 19:04:29                   │
│                                      │
│  Текущее сохранение будет            │
│  перезаписано.                       │
│                                      │
│  [ ОТМЕНА ]    [ ВОССТАНОВИТЬ ]      │
│                                      │
└──────────────────────────────────────┘
```

### Реализация

- Оверлей: `.modal-overlay` (fixed, `rgba(0,0,0,0.85)`, z-index 1000)
- Скрытие: класс `.hidden` (`display: none`)
- Контейнер: `.modal-box` (фон `#141c28`, max-width 400px)
- Кнопки: `.btn-cancel` + `.btn-restore`
- JavaScript:
  - `showModal()` — убирает `.hidden`
  - `hideModal()` — добавляет `.hidden`
  - `doRestore()` — `fetch('/restore/{id}', {method:'POST'})` → JSON → redirect на `/`
  - `Escape` закрывает модал

---

## Вспомогательные функции (Go)

Файл: `app/appList.go`

| Функция | Описание |
|---------|----------|
| `backupTime(path)` | Извлекает `HH:MM` из имени папки бэкапа |
| `backupDate(path)` | Извлекает `YYYY-MM-DD` из имени папки |
| `formatDateRu(iso)` | `"2026-02-19"` → `"19 февраля 2026"` |
| `condClass(s)` | `> 60` → `cond-good`, `> 30` → `cond-warn`, иначе `cond-danger` |
| `condWidth(s)` | Значение → `"52%"` (clamped 0-100) |
| `pluralBackups(n)` | Русское склонение: "1 бэкап", "2 бэкапа", "5 бэкапов" |

---

## CSS-файл

Файл: `app/web/style.css`

### Секции

| Секция | Классы |
|--------|--------|
| Reset | `*`, `body`, `a` |
| Header | `.header`, `.header h1` |
| Hero Card | `.hero`, `.hero-img`, `.hero-info`, `.hero-title`, `.hero-subtitle`, `.hero-stats`, `.stat-row`, `.stat-label`, `.stat-value`, `.hero-footer`, `.backup-time` |
| Condition Bar | `.condition-bar`, `.condition-fill`, `.cond-good`, `.cond-warn`, `.cond-danger` |
| Buttons | `.btn-restore`, `.btn-cancel` |
| Save Section | `.save-section`, `.save-header`, `.save-name`, `.save-meta`, `.save-subtitle`, `.save-body` |
| Date Section | `.date-section`, `.date-separator`, `.date-text`, `.date-count` |
| Card Grid | `.card-grid` |
| Card | `.card`, `.card-img`, `.card-body`, `.card-time`, `.card-region`, `.card-survived`, `.card-condition`, `.card-condition-text` |
| Detail Page | `.detail-page`, `.back-link`, `.detail-content`, `.detail-img`, `.detail-info`, `.detail-title`, `.detail-mode`, `.detail-stats`, `.detail-stat`, `.detail-stat-label`, `.detail-stat-value`, `.detail-meta` |
| Modal | `.modal-overlay`, `.modal-box`, `.modal-warning`, `.modal-buttons` |
| Utilities | `.text-error` |
| Responsive | `@media (max-width: 768px)`, `@media (max-width: 480px)` |
