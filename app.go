package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

var (
	IgnFileList = []string{
		".DS_Store",
	}
)

func StringOnList(val string, list []string) bool {
	for _, itm := range list {
		if strings.Contains(itm, val) {
			return true
		}
	}

	return false
}

func statFile(Data []byte, MIME string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", MIME)
		w.Write(Data)
	})
}

func appSrv() error {

	app.Route("/", func() app.Composer { return &ListPage{} })
	//app.Route("/list/", func() app.Composer { return &ListPage{} })
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "LG Backup",
		Description: "Automatic backup service",
		Styles: []string{
			"/web/style.css",
		},
		Icon: app.Icon{
			Default: "/web/favicon.png",
		},
		//Scripts:     []string{"/js/main.js"}
		//Resources: app.LocalDir("/"),
	})

	http.Handle("/web/favicon.png", statFile(cFavIconPNG, "image/png"))
	http.Handle("/web/style.css", statFile(cStyleCSS, "text/css"))
	http.Handle("/web/app.wasm", statFile(cAppWASM, "application/wasm"))
	http.HandleFunc("/restore/", DoRestoreHandler)

	go func() {
		log.Println("Listening on " + cHost)

		if err := http.ListenAndServe(cHost, nil); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
