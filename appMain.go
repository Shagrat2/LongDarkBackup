package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type MainPage struct {
	app.Compo
}

func (h *MainPage) Render() app.UI {

	list := make([]app.UI, 0)

	err := filepath.WalkDir(BackupFolder, func(path string, d os.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}

		// Skip file
		if !strings.HasPrefix(d.Name(), "sandbox") {
			return nil
		}

		fPath := strings.Trim(filepath.Dir(path[len(BackupFolder):]), "/\\")

		info, img, _ := loadData(fPath)

		list = append(list,
			app.Details().Body(
				app.Summary().Body(
					app.Span().Text(fPath), app.Span().Text(info.DisplayName), app.Br(),
					app.Span().Text(info.XPMode), app.Span().Text(info.Region), app.Br(),
					app.Span().Text(info.HoursSurvived),
				),
				app.P().Body(
					app.Span().Text("Condition: "+info.Condition), app.Br(),
					app.Span().Text("World explored: "+info.WorldExplored), app.Br(),
					app.Span().Text("Game mode: "+info.GameMode), app.Br(),
					app.Img().Src("data:image/png;base64,"+base64.StdEncoding.EncodeToString(img)),
					app.Br(),
					app.A().Text("Restore").Href("/restore/"+fPath),
				),
			),
		)

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

	return app.Div().
		Class("container").
		Body(
			app.H1().Text("LongDark Backup"),
			app.Div().Body(
				list...,
			),
		)
}
