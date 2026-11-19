package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type SettingsData struct {
	BackupDir string
	Saved     bool
	Error     string
}

var settingsTmpl = newPageTemplate("settings", `<!DOCTYPE html>
<html lang="{{Lang}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{T "settings"}} — LongDark Backup</title>
<link rel="stylesheet" href="/web/style.css">
<link rel="icon" href="/web/favicon.png">
</head>
<body>
<div class="container">
  <div class="detail-page">
    <div class="detail-top">
      <a href="/" class="back-link">&#8592; {{T "back"}}</a>
      {{template "lang_switch" .}}
    </div>

    <div class="detail-title">{{T "settings"}}</div>

    {{if .Saved}}<div class="settings-msg settings-ok">{{T "settings_saved"}}</div>{{end}}
    {{if .Error}}<div class="settings-msg settings-err">{{.Error}}</div>{{end}}

    <form method="post" action="/settings" class="settings-form">
      <label class="settings-label" for="backup_dir">{{T "backup_dir"}}</label>
      <input class="settings-input" type="text" id="backup_dir" name="backup_dir" value="{{.BackupDir}}" spellcheck="false" autocomplete="off">
      <p class="settings-hint">{{T "backup_dir_hint"}}</p>
      <button class="btn-restore" type="submit">{{T "save"}}</button>
    </form>
  </div>
</div>
</body>
</html>`)

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && !isSameOrigin(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var data SettingsData

	if r.Method == http.MethodPost {
		dir := strings.TrimSpace(r.FormValue("backup_dir"))
		if dir != "" {
			dir = filepath.Clean(dir)
		}
		data.BackupDir = dir

		switch {
		case dir == "":
			data.Error = T(MsgPathRequired)
		case !filepath.IsAbs(dir):
			data.Error = T(MsgPathRelative)
		case filepath.Dir(dir) == dir:
			data.Error = T(MsgPathRoot)
		default:
			if key, err := prepareFolder(dir); err != nil {
				data.Error = errorText(key, err)
			} else if key, err := applyBackupFolder(dir); err != nil {
				data.Error = errorText(key, err)
			} else {
				data.Saved = true
			}
		}
	} else {
		data.BackupDir = GetBackupFolder()
	}

	renderHTML(w, settingsTmpl, data)
}

// prepareFolder creates the directory if needed and verifies it is writable.
// It returns the translation key describing the failed step.
func prepareFolder(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return MsgDirCreateErr, fmt.Errorf("%s: %w", dir, err)
	}

	f, err := os.CreateTemp(dir, ".ldb-check-*")
	if err != nil {
		return MsgDirNotWritable, fmt.Errorf("%s: %w", dir, err)
	}
	name := f.Name()
	if err := f.Close(); err != nil {
		os.Remove(name)
		return MsgDirNotWritable, fmt.Errorf("%s: %w", dir, err)
	}
	if err := os.Remove(name); err != nil {
		return MsgDirNotWritable, fmt.Errorf("%s: %w", dir, err)
	}

	return "", nil
}

// applyBackupFolder saves the directory to the config and applies it at runtime.
func applyBackupFolder(dir string) (string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		log.Println("Load config error:", err)
		if rerr := BackupBrokenConfig(); rerr != nil {
			return MsgConfigReadErr, err
		}
		cfg = Config{}
	}

	cfg.BackupDir = dir
	if err := SaveConfig(cfg); err != nil {
		return MsgConfigSaveErr, err
	}

	SetBackupFolder(dir)
	return "", nil
}

func errorText(key string, err error) string {
	return T(key) + ": " + err.Error()
}
