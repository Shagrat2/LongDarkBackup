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

func statFile(Data []byte, MIME string, cache string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", MIME)
		w.Header().Set("Cache-Control", cache)
		w.Write(Data)
	})
}

func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		next.ServeHTTP(w, r)
	})
}

func appSrv() error {

	app.Route("/", func() app.Composer { return &ListPage{} })
	//app.Route("/list/", func() app.Composer { return &ListPage{} })
	app.RunWhenOnBrowser()

	http.Handle("/", noCache(&app.Handler{
		Name:        "LG Backup",
		Description: "Automatic backup service",
		Styles: []string{
			"/web/style.css",
		},
		Icon: app.Icon{
			Default: "/web/favicon.png",
		},
	}))

	http.Handle("/web/favicon.png", statFile(cFavIconPNG, "image/png", "max-age=86400"))
	http.Handle("/web/style.css", statFile(cStyleCSS, "text/css", "no-cache, no-store, must-revalidate"))
	http.Handle("/web/app.wasm", statFile(cAppWASM, "application/wasm", "no-cache, no-store, must-revalidate"))
	http.HandleFunc("/img/", ImgHandler)
	http.HandleFunc("/save/", SaveDetailHandler)
	http.HandleFunc("/restore/", DoRestoreHandler)

	go func() {
		log.Println("Listening on " + cHost)

		if err := http.ListenAndServe(cHost, nil); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
