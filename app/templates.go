package main

import (
	"html/template"
	"log"
	"net/http"
)

const langSwitchSrc = `<div class="lang-switch">
        <a class="lang-link{{if eq Lang "ru"}} active{{end}}" href="/set-lang?lang=ru" onclick="localStorage.setItem('lang','ru')">RU</a>
        <span class="lang-sep">/</span>
        <a class="lang-link{{if eq Lang "en"}} active{{end}}" href="/set-lang?lang=en" onclick="localStorage.setItem('lang','en')">EN</a>
      </div>`

// newPageTemplate builds a page template with the shared language switcher
// available as {{template "lang_switch" .}}.
func newPageTemplate(name, body string) *template.Template {
	t := template.New("lang_switch").Funcs(template.FuncMap{
		"T":    T,
		"Lang": func() string { return string(GetLang()) },
	})
	template.Must(t.Parse(langSwitchSrc))
	template.Must(t.New(name).Parse(body))
	return t.Lookup(name)
}

// renderHTML writes an HTML page with the caching policy required by the UI.
func renderHTML(w http.ResponseWriter, tmpl *template.Template, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Template error:", err)
	}
}
