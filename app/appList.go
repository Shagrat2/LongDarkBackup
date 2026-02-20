package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type ListPage struct {
	app.Compo
}

type TItem struct {
	Info       cacheItem
	Path       string
	BackupTime time.Time
}

type TDateGroup struct {
	Date      string
	DateLabel string
	Items     []*TItem
}

type TSaveGroup struct {
	Name  string
	Items []*TItem
	Days  []TDateGroup
}

// ── Helpers ──

var ruMonths = []string{
	"", "января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

var enMonths = []string{
	"", "January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

func backupTime(path string) string {
	if len(path) >= 16 {
		return path[11:13] + ":" + path[14:16]
	}
	return ""
}

func backupDate(path string) string {
	if len(path) >= 10 {
		return path[:10]
	}
	return ""
}

func formatDate(iso string) string {
	t, err := time.Parse("2006-01-02", iso)
	if err != nil {
		return iso
	}
	if GetLang() == LangEN {
		return fmt.Sprintf("%s %d, %d", enMonths[t.Month()], t.Day(), t.Year())
	}
	return fmt.Sprintf("%d %s %d", t.Day(), ruMonths[t.Month()], t.Year())
}

func condClass(s string) string {
	v, _ := strconv.ParseFloat(s, 64)
	if v > 60 {
		return "cond-good"
	}
	if v > 30 {
		return "cond-warn"
	}
	return "cond-danger"
}

func condWidth(s string) string {
	v, _ := strconv.ParseFloat(s, 64)
	if v < 0 {
		v = 0
	}
	if v > 100 {
		v = 100
	}
	return fmt.Sprintf("%.0f%%", v)
}

func pluralBackups(n int) string {
	if GetLang() == LangEN {
		if n == 1 {
			return "1 backup"
		}
		return fmt.Sprintf("%d backups", n)
	}
	if n%10 == 1 && n%100 != 11 {
		return fmt.Sprintf("%d бэкап", n)
	}
	if n%10 >= 2 && n%10 <= 4 && (n%100 < 10 || n%100 >= 20) {
		return fmt.Sprintf("%d бэкапа", n)
	}
	return fmt.Sprintf("%d бэкапов", n)
}

// ── Render ──

func (h *ListPage) Render() app.UI {

	// In browser (WASM) context, sync language from localStorage
	if ls := app.Window().Get("localStorage"); ls.Truthy() {
		if stored := ls.Call("getItem", "lang"); stored.Truthy() {
			SetLang(DetectLang(stored.String(), ""))
		} else {
			// First visit: read cookie set by middleware, save to localStorage
			for _, c := range strings.Split(app.Window().Get("document").Get("cookie").String(), ";") {
				c = strings.TrimSpace(c)
				if strings.HasPrefix(c, "lang=") {
					ls.Call("setItem", "lang", c[5:])
					SetLang(DetectLang(c[5:], ""))
					break
				}
			}
		}
	}

	// Collect all backup items
	var allItems []*TItem

	filepath.WalkDir(BackupFolder, func(path string, d os.DirEntry, err error) error {
		if d == nil || !d.IsDir() {
			return nil
		}
		if !FileExists(path, "sandbox") || !FileExists(path, "profile") {
			return nil
		}

		fPath := filepath.Base(path)
		info, _, _ := loadData(fPath)
		if info.DisplayName == "" {
			return nil
		}

		var bt time.Time
		if fi, err := os.Stat(filepath.Join(BackupFolder, fPath, "info.json")); err == nil {
			bt = fi.ModTime()
		}

		allItems = append(allItems, &TItem{
			Info:       info,
			Path:       fPath,
			BackupTime: bt,
		})
		return nil
	})

	// Sort newest first
	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].BackupTime.After(allItems[j].BackupTime)
	})

	// Group by save name
	saveMap := make(map[string]*TSaveGroup)
	var saveOrder []string

	for _, itm := range allItems {
		name := itm.Info.DisplayName
		sg, ok := saveMap[name]
		if !ok {
			sg = &TSaveGroup{Name: name}
			saveMap[name] = sg
			saveOrder = append(saveOrder, name)
		}
		sg.Items = append(sg.Items, itm)
	}

	// Sort save groups by most recent backup
	sort.Slice(saveOrder, func(i, j int) bool {
		return saveMap[saveOrder[i]].Items[0].BackupTime.After(saveMap[saveOrder[j]].Items[0].BackupTime)
	})

	// Build date groups within each save
	var saves []*TSaveGroup
	for _, name := range saveOrder {
		sg := saveMap[name]

		for _, itm := range sg.Items {
			d := backupDate(itm.Path)
			if len(sg.Days) == 0 || sg.Days[len(sg.Days)-1].Date != d {
				sg.Days = append(sg.Days, TDateGroup{
					Date:      d,
					DateLabel: formatDate(d),
				})
			}
			sg.Days[len(sg.Days)-1].Items = append(sg.Days[len(sg.Days)-1].Items, itm)
		}

		saves = append(saves, sg)
	}

	// Language switcher
	lang := GetLang()
	var ruClass, enClass string
	if lang == LangRU {
		ruClass = "lang-link active"
		enClass = "lang-link"
	} else {
		ruClass = "lang-link"
		enClass = "lang-link active"
	}

	// Build UI
	content := []app.UI{
		app.Div().Class("header").Body(
			app.H1().Text(T(MsgAppTitle)),
			app.Div().Class("lang-switch").Body(
				app.A().Class(ruClass).Href("/set-lang?lang=ru").
					Attr("onclick", "localStorage.setItem('lang','ru')").Text("RU"),
				app.Span().Class("lang-sep").Text("/"),
				app.A().Class(enClass).Href("/set-lang?lang=en").
					Attr("onclick", "localStorage.setItem('lang','en')").Text("EN"),
			),
		),
	}

	onlyOne := len(saves) == 1

	for _, sg := range saves {
		sectionContent := []app.UI{}

		// Date groups with card grids (collapsed by default)
		for _, dg := range sg.Days {
			cards := make([]app.UI, 0, len(dg.Items))
			for _, itm := range dg.Items {
				cards = append(cards, renderCard(itm))
			}

			sectionContent = append(sectionContent,
				app.Details().Class("date-section").Body(
					app.Summary().Class("date-separator").Body(
						app.Span().Class("date-text").Text(dg.DateLabel),
						app.Span().Class("date-count").Text(pluralBackups(len(dg.Items))),
					),
					app.Div().Class("card-grid").Body(cards...),
				),
			)
		}

		// Subtitle for the section header
		latest := sg.Items[0]
		subtitle := latest.Info.XPMode
		if latest.Info.Region != "" {
			subtitle += " · " + latest.Info.Region
		}

		details := app.Details().Class("save-section").Body(
			app.Summary().Class("save-header").Body(
				app.Div().Body(
					app.Span().Class("save-name").Text(sg.Name),
					app.Span().Class("save-subtitle").Text(subtitle),
				),
				app.Span().Class("save-meta").Text(pluralBackups(len(sg.Items))),
			),
			app.Div().Class("save-body").Body(sectionContent...),
		)

		// Auto-open if only one save
		if onlyOne {
			details = details.Attr("open", "")
		}

		content = append(content, details)
	}

	if len(saves) == 0 {
		content = append(content,
			app.P().Class("text-error").Text(T(MsgNoSaves)),
		)
	}

	return app.Div().Class("container").Attr("data-lang", string(GetLang())).Body(content...)
}

func renderCard(itm *TItem) app.UI {
	return app.A().Class("card").Href("/save/"+itm.Path).Body(
		app.Img().Class("card-img").Src("/img/"+itm.Path).Alt("Screenshot").Attr("loading", "lazy"),
		app.Div().Class("card-body").Body(
			app.Div().Class("card-time").Text(backupTime(itm.Path)),
			app.Div().Class("card-region").Text(itm.Info.Region),
			app.Div().Class("card-survived").Text(itm.Info.HoursSurvived),
			app.Div().Class("card-condition").Body(
				app.Div().Class("condition-bar "+condClass(itm.Info.Condition)).Body(
					app.Div().Class("condition-fill").
						Style("width", condWidth(itm.Info.Condition)),
				),
				app.Span().Class("card-condition-text").Text(itm.Info.Condition+"%"),
			),
		),
	)
}
