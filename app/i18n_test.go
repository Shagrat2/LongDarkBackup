package main

import "testing"

func TestTranslationsParity(t *testing.T) {
	ru, en := translations[LangRU], translations[LangEN]

	for k := range ru {
		if _, ok := en[k]; !ok {
			t.Errorf("missing EN translation for %q", k)
		}
	}
	for k := range en {
		if _, ok := ru[k]; !ok {
			t.Errorf("missing RU translation for %q", k)
		}
	}
}
