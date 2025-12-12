package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type ListPage struct {
	app.Compo
}

type TSaves []*TSave

type TSave struct {
	Title string
	Days  TDays
}

type TDays []*TDay

type TDay struct {
	Date  string
	Items TItems
}

type TItems []*TItem

type TItem struct {
	Info cacheItem
	Path string
	Img  []byte
}

// --- Items
func (a TItems) Len() int           { return len(a) }
func (a TItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TItems) Less(i, j int) bool { return a[i].Info.Timestamp < a[j].Info.Timestamp }

func (a TItems) Find(val string) *TItem {
	for _, itm := range a {
		if itm.Info.Timestamp == val {
			return itm
		}
	}

	return nil
}

// --- Days
func (a TDay) Sort() {
	sort.Sort(a.Items)
}

func (a TDays) Len() int           { return len(a) }
func (a TDays) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TDays) Less(i, j int) bool { return a[i].Date < a[j].Date }

func (a TDays) Find(val string) *TDay {
	for _, itm := range a {
		if itm.Date == val {
			return itm
		}
	}

	return nil
}

// --- Saves
func (a TSave) Sort() {
	sort.Sort(a.Days)
}

func (a TSaves) Len() int           { return len(a) }
func (a TSaves) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TSaves) Less(i, j int) bool { return a[i].Title < a[j].Title }

func (a TSaves) Find(val string) *TSave {
	for _, itm := range a {
		if itm.Title == val {
			return itm
		}
	}

	return nil
}

func (a TSaves) Sort() {

	sort.Sort(a)

	for _, itm := range a {
		itm.Sort()
	}
}

// ListPage
func (h *ListPage) Render() app.UI {

	var fSaves TSaves

	err := filepath.WalkDir(BackupFolder, func(path string, d os.DirEntry, err error) error {
		if d == nil || !d.IsDir() {
			return nil
		}

		// Check file
		fHasSandbox := FileExists(path, "sandbox")
		fHasProfile := FileExists(path, "profile")

		if !fHasSandbox || !fHasProfile {
			return nil
		}

		fPath := filepath.Base(path)

		// log.Println("Parse:", fPath)

		info, img, _ := loadData(fPath)

		fArr := strings.Split(info.Timestamp, " ")
		if len(fArr) < 1 {
			return nil
		}
		fDate := fArr[0]

		// Get save
		fSave := fSaves.Find(info.DisplayName)
		if fSave == nil {
			fSave = &TSave{
				Title: info.DisplayName,
			}
			fSaves = append(fSaves, fSave)
		} // fSave

		// Get day
		fDay := fSave.Days.Find(fDate)
		if fDay == nil {
			fDay = &TDay{
				Date: fDate,
			}
			fSave.Days = append(fSave.Days, fDay)
		} // fDay

		// Add save
		fDay.Items = append(fDay.Items, &TItem{
			Info: info,
			Path: fPath,
			Img:  img,
		})

		return nil
	})
	if err != nil {

		return app.Div().
			Class("container").
			Body(
				app.P().
					Class("text-error").
					Text("Error reading folder"),
			)
	}

	fSaves.Sort()

	list := make([]app.UI, 0)
	for _, fSave := range fSaves {

		iDays := make([]app.UI, 0)

		for _, fDay := range fSave.Days {

			iItems := make([]app.UI, 0)
			for _, fItm := range fDay.Items {

				// Items
				iItems = append(iItems, app.Details().Body(
					app.Summary().Body(
						// app.Span().Text(fPath), app.Span().Text(info.DisplayName), app.Br(),
						app.Span().Text(fItm.Info.XPMode), app.Span().Text(fItm.Info.Region), app.Br(),
						app.Span().Text(fItm.Info.HoursSurvived),
					),
					app.P().Body(
						app.Span().Text("Condition: "+fItm.Info.Condition), app.Br(),
						app.Span().Text("World explored: "+fItm.Info.WorldExplored), app.Br(),
						app.Span().Text("Game mode: "+fItm.Info.GameMode), app.Br(),
						app.Img().Src("data:image/jpeg;base64,"+base64.StdEncoding.EncodeToString(fItm.Img)),
						app.Br(),
						app.A().Text("Restore").Href("/restore/"+fItm.Path),
					),
				))
			}

			// Days
			iDays = append(iDays, app.Details().Body(
				app.Summary().Body(
					app.Span().Text(fDay.Date),
				),
				app.P().Body(
					iItems...,
				),
			))

		}

		list = append(list, app.Details().Body(
			app.Summary().Body(
				app.Span().Text(fSave.Title),
			),
			app.P().Body(
				iDays...,
			),
		))
	}

	return app.Div().
		Class("container").
		Body(
			app.H1().Text("LongDark Backup"),
			app.Div().Body(
				list...,
			),
		)
}
