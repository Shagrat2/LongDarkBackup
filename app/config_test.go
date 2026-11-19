package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupConfigDir(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)
}

func postSettings(t *testing.T, dir string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	form := url.Values{"backup_dir": {dir}}
	req := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	SettingsHandler(w, req)
	return w
}

func TestResolveFolder(t *testing.T) {
	cases := []struct {
		configValue, envValue, defaultValue, want string
	}{
		{"cfg", "env", "def", "cfg"},
		{"", "env", "def", "env"},
		{"", "", "def", "def"},
	}
	for _, c := range cases {
		if got := ResolveFolder(c.configValue, c.envValue, c.defaultValue); got != c.want {
			t.Errorf("ResolveFolder(%q, %q, %q) = %q, want %q",
				c.configValue, c.envValue, c.defaultValue, got, c.want)
		}
	}
}

func TestConfigRoundTrip(t *testing.T) {
	setupConfigDir(t)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BackupDir != "" {
		t.Fatalf("expected empty config, got %q", cfg.BackupDir)
	}

	dir := filepath.Join(t.TempDir(), "backups")
	if err := SaveConfig(Config{DataDir: "/data", BackupDir: dir}); err != nil {
		t.Fatal(err)
	}

	cfg, err = LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BackupDir != dir || cfg.DataDir != "/data" {
		t.Fatalf("unexpected config: %+v", cfg)
	}

	path, _ := ConfigPath()
	tmp, _ := filepath.Glob(filepath.Join(filepath.Dir(path), "config-*.tmp"))
	if len(tmp) != 0 {
		t.Fatalf("temp files left behind: %v", tmp)
	}
}

func TestSettingsHandler(t *testing.T) {
	setupConfigDir(t)
	SetLang(LangEN)
	SetBackupFolder("/old")

	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	w := httptest.NewRecorder()
	SettingsHandler(w, req)
	if !strings.Contains(w.Body.String(), "/old") {
		t.Fatalf("GET body missing current dir: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "lang-switch") {
		t.Fatalf("GET body missing language switcher: %s", w.Body.String())
	}

	dir := filepath.Join(t.TempDir(), "My Backups")
	w = postSettings(t, dir+string(os.PathSeparator), nil)
	if !strings.Contains(w.Body.String(), T(MsgSettingsSaved)) {
		t.Fatalf("save failed: %s", w.Body.String())
	}
	if GetBackupFolder() != dir {
		t.Fatalf("runtime folder = %q, want %q", GetBackupFolder(), dir)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("backup dir not created: %v", err)
	}
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BackupDir != dir {
		t.Fatalf("config backup dir = %q, want %q", cfg.BackupDir, dir)
	}
}

func TestSettingsRejectsInvalidPaths(t *testing.T) {
	setupConfigDir(t)
	SetLang(LangEN)
	SetBackupFolder("/old")

	for _, dir := range []string{"", "relative/path", string(os.PathSeparator)} {
		w := postSettings(t, dir, nil)
		if !strings.Contains(w.Body.String(), "settings-err") {
			t.Errorf("path %q: expected error, got: %s", dir, w.Body.String())
		}
		if GetBackupFolder() != "/old" {
			t.Errorf("path %q: folder changed to %q", dir, GetBackupFolder())
		}
	}
}

func TestSettingsRejectsCrossOrigin(t *testing.T) {
	setupConfigDir(t)
	SetLang(LangEN)
	SetBackupFolder("/old")

	dir := filepath.Join(t.TempDir(), "backups")
	w := postSettings(t, dir, map[string]string{"Origin": "http://evil.example"})
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
	if GetBackupFolder() != "/old" {
		t.Fatalf("folder changed to %q", GetBackupFolder())
	}
}

func TestSettingsRepairsBrokenConfig(t *testing.T) {
	setupConfigDir(t)
	SetLang(LangEN)
	SetBackupFolder("/old")

	path, err := ConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{broken"), 0644); err != nil {
		t.Fatal(err)
	}

	dir := filepath.Join(t.TempDir(), "backups")
	w := postSettings(t, dir, nil)
	if !strings.Contains(w.Body.String(), T(MsgSettingsSaved)) {
		t.Fatalf("save failed: %s", w.Body.String())
	}
	if _, err := os.Stat(path + ".bak"); err != nil {
		t.Fatalf("broken config was not backed up: %v", err)
	}
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BackupDir != dir {
		t.Fatalf("config backup dir = %q, want %q", cfg.BackupDir, dir)
	}
}

func TestIsSameOrigin(t *testing.T) {
	cases := []struct {
		name    string
		headers map[string]string
		want    bool
	}{
		{"no headers", nil, true},
		{"same origin", map[string]string{"Origin": "http://example.com"}, true},
		{"cross origin", map[string]string{"Origin": "http://evil.example"}, false},
		{"bad origin", map[string]string{"Origin": "://bad"}, false},
		{"same site", map[string]string{"Sec-Fetch-Site": "same-origin"}, true},
		{"cross site", map[string]string{"Sec-Fetch-Site": "cross-site"}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/settings", nil)
			for k, v := range c.headers {
				req.Header.Set(k, v)
			}
			if got := isSameOrigin(req); got != c.want {
				t.Fatalf("isSameOrigin() = %v, want %v", got, c.want)
			}
		})
	}
}
