package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type DetailData struct {
	Info          cacheItem
	Global        GlobalData
	BackupDateFmt string
	BackupTimeFmt string
	CondClass     string
	CondWidth     string
	ID            string
	Lang          string
}

var detailTmpl = template.Must(template.New("detail").Funcs(template.FuncMap{
	"T": T,
}).Parse(`<!DOCTYPE html>
<html lang="{{.Lang}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Info.Region}} — LongDark Backup</title>
<link rel="stylesheet" href="/web/style.css">
<link rel="icon" href="/web/favicon.png">
</head>
<body>
<div class="container">
  <div class="detail-page">
    <div class="detail-top">
      <a href="#" class="back-link" onclick="history.back();return false">&#8592; {{T "back"}}</a>
      <div class="lang-switch">
        <a class="lang-link{{if eq .Lang "ru"}} active{{end}}" href="/set-lang?lang=ru" onclick="localStorage.setItem('lang','ru')">RU</a>
        <span class="lang-sep">/</span>
        <a class="lang-link{{if eq .Lang "en"}} active{{end}}" href="/set-lang?lang=en" onclick="localStorage.setItem('lang','en')">EN</a>
      </div>
    </div>

    <div class="detail-content">
      <img class="detail-img" src="/img/{{.ID}}" alt="Screenshot">
      <div class="detail-info">
        <div class="detail-title">{{.Info.Region}}</div>
        <div class="detail-mode">{{.Info.XPMode}} &middot; {{.Info.GameMode}}</div>

        <div class="detail-stats">
          <div class="detail-stat">
            <span class="detail-stat-label">{{T "survival"}}</span>
            <span class="detail-stat-value">{{.Info.HoursSurvived}}</span>
          </div>
          <div class="detail-stat">
            <span class="detail-stat-label">{{T "condition"}}</span>
            <span class="detail-stat-value">
              {{.Info.Condition}}%
              <span class="condition-bar {{.CondClass}}"><span class="condition-fill" style="width:{{.CondWidth}}"></span></span>
            </span>
          </div>
          {{if .Global.Hunger}}<div class="detail-stat">
            <span class="detail-stat-label">{{T "hunger"}}</span>
            <span class="detail-stat-value">{{.Global.Hunger}}</span>
          </div>{{end}}
          {{if .Global.Thirst}}<div class="detail-stat">
            <span class="detail-stat-label">{{T "thirst"}}</span>
            <span class="detail-stat-value">{{.Global.Thirst}}</span>
          </div>{{end}}
          {{if .Global.Fatigue}}<div class="detail-stat">
            <span class="detail-stat-label">{{T "fatigue"}}</span>
            <span class="detail-stat-value">{{.Global.Fatigue}}</span>
          </div>{{end}}
          {{if .Global.Freezing}}<div class="detail-stat">
            <span class="detail-stat-label">{{T "freezing"}}</span>
            <span class="detail-stat-value">{{.Global.Freezing}}</span>
          </div>{{end}}
          <div class="detail-stat">
            <span class="detail-stat-label">{{T "world_explored"}}</span>
            <span class="detail-stat-value">{{.Info.WorldExplored}}%</span>
          </div>
          <div class="detail-stat">
            <span class="detail-stat-label">{{T "character"}}</span>
            <span class="detail-stat-value">{{.Info.Persona}}</span>
          </div>
        </div>
        {{if .Global.Afflictions}}
        <div class="detail-afflictions">
          {{range .Global.Afflictions}}<span class="affliction-tag">{{.}}</span>{{end}}
        </div>
        {{end}}

        <div class="detail-meta">
          {{T "backup_label"}}: {{.BackupDateFmt}}, {{.BackupTimeFmt}}<br>
          {{T "save_label"}}: {{.Info.DisplayName}} &middot; {{.Info.Timestamp}}
        </div>

        <button class="btn-restore" onclick="showModal()">&#8635; {{T "restore"}}</button>
      </div>
    </div>
  </div>
</div>

<div id="modal" class="modal-overlay hidden">
  <div class="modal-box">
    <h2>{{T "restore_confirm"}}</h2>
    <p>{{.Info.Region}} &middot; {{.Info.HoursSurvived}}</p>
    <p>{{T "backup_from"}} {{.BackupTimeFmt}}</p>
    <p class="modal-warning">{{T "overwrite_warn"}}</p>
    <div class="modal-buttons">
      <button class="btn-cancel" onclick="hideModal()">{{T "cancel"}}</button>
      <button class="btn-restore" onclick="doRestore()">{{T "restore"}}</button>
    </div>
  </div>
</div>

<script>
function showModal(){document.getElementById('modal').classList.remove('hidden')}
function hideModal(){document.getElementById('modal').classList.add('hidden')}
function doRestore(){
  fetch('/restore/{{.ID}}',{method:'POST'})
  .then(function(r){return r.json()})
  .then(function(d){
    if(d.ok){window.location='/'}
    else{alert(d.error);hideModal()}
  })
  .catch(function(e){alert(e);hideModal()})
}
document.addEventListener('keydown',function(e){if(e.key==='Escape')hideModal()})
</script>
</body>
</html>`))

func ImgHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)

	_, img, err := loadData(id)
	if err != nil || len(img) == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Write(img)
}

func SaveDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)

	info, _, err := loadData(id)
	if err != nil {
		http.Error(w, "Save not found", http.StatusNotFound)
		return
	}

	bDate := backupDate(id)
	bTime := ""
	if len(id) >= 19 {
		bTime = id[11:13] + ":" + id[14:16] + ":" + id[17:19]
	}

	data := DetailData{
		Info:          info,
		Global:        loadGlobalData(id),
		BackupDateFmt: formatDate(bDate),
		BackupTimeFmt: bTime,
		CondClass:     condClass(info.Condition),
		CondWidth:     condWidth(info.Condition),
		ID:            id,
		Lang:          string(GetLang()),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	detailTmpl.Execute(w, data)
}
