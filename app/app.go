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

func noETag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("If-None-Match")
		next.ServeHTTP(w, r)
	})
}

func langMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("lang")
		cookieVal := ""
		if cookie != nil {
			cookieVal = cookie.Value
		}
		lang := DetectLang(cookieVal, r.Header.Get("Accept-Language"))
		SetLang(lang)
		if cookie == nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "lang",
				Value:  string(lang),
				Path:   "/",
				MaxAge: 365 * 24 * 3600,
			})
		}
		next.ServeHTTP(w, r)
	})
}

func setLangHandler(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang != "ru" && lang != "en" {
		lang = "en"
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "lang",
		Value:  lang,
		Path:   "/",
		MaxAge: 365 * 24 * 3600,
	})
	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusFound)
}

func appSrv() error {

	app.Route("/", func() app.Composer { return &ListPage{} })
	//app.Route("/list/", func() app.Composer { return &ListPage{} })
	app.RunWhenOnBrowser()

	appHandler := &app.Handler{
		Name:        "LG Backup",
		Description: "Automatic backup service",
		Styles: []string{
			"/web/style.css",
		},
		Icon: app.Icon{
			Default: "/web/favicon.png",
		},
		RawHeaders: []string{
			`<style>#app-wasm-loader{display:none!important}</style>`,
			`<script>window.addEventListener('pageshow',function(e){if(e.persisted){var l=localStorage.getItem('lang');var el=document.querySelector('[data-lang]');if(l&&el&&el.dataset.lang!==l)location.reload()}})</script>`,
		},
	}

	http.Handle("/", langMiddleware(noETag(appHandler)))

	http.Handle("/web/favicon.png", statFile(cFavIconPNG, "image/png", "max-age=86400"))
	http.Handle("/web/style.css", statFile(cStyleCSS, "text/css", "no-cache, no-store, must-revalidate"))
	http.HandleFunc("/web/app.wasm", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	http.HandleFunc("/img/", ImgHandler)
	http.Handle("/save/", langMiddleware(http.HandlerFunc(SaveDetailHandler)))
	http.HandleFunc("/restore/", DoRestoreHandler)
	http.HandleFunc("/set-lang", setLangHandler)

	go func() {
		log.Println("Listening on " + cHost)

		if err := http.ListenAndServe(cHost, nil); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
