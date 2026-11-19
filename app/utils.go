package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

func FileExists(Dir string, FilaPref string) bool {

	items, err := os.ReadDir(Dir)
	if err != nil {
		return false
	}

	for _, itm := range items {
		if strings.HasPrefix(itm.Name(), FilaPref) {
			return true
		}
	}

	return false
}

// isSameOrigin reports whether a state-changing request comes from this app.
// Requests without Origin/Sec-Fetch-Site headers (curl, older clients) are allowed.
func isSameOrigin(r *http.Request) bool {
	if origin := r.Header.Get("Origin"); origin != "" {
		u, err := url.Parse(origin)
		if err != nil {
			return false
		}
		return strings.EqualFold(u.Host, r.Host)
	}

	switch r.Header.Get("Sec-Fetch-Site") {
	case "", "same-origin", "none":
		return true
	}

	return false
}
