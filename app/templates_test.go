package main

import (
	"strings"
	"testing"
)

func TestDetailTemplateRenders(t *testing.T) {
	SetLang(LangEN)

	var b strings.Builder
	data := DetailData{
		Info: cacheItem{Region: "Mystery Lake"},
		ID:   "2026-01-01-00-00-00",
	}
	if err := detailTmpl.Execute(&b, data); err != nil {
		t.Fatal(err)
	}

	out := b.String()
	for _, want := range []string{"Mystery Lake", "lang-switch", "/set-lang?lang=ru", "/set-lang?lang=en"} {
		if !strings.Contains(out, want) {
			t.Fatalf("rendered page missing %q", want)
		}
	}
}
