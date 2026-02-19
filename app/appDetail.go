package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type DetailData struct {
	Info          cacheItem
	BackupDateFmt string
	BackupTimeFmt string
	CondClass     string
	CondWidth     string
	ID            string
}

var detailTmpl = template.Must(template.New("detail").Parse(`<!DOCTYPE html>
<html lang="ru">
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
    <a href="/" class="back-link">&#8592; НАЗАД</a>

    <div class="detail-content">
      <img class="detail-img" src="/img/{{.ID}}" alt="Screenshot">
      <div class="detail-info">
        <div class="detail-title">{{.Info.Region}}</div>
        <div class="detail-mode">{{.Info.XPMode}} &middot; {{.Info.GameMode}}</div>

        <div class="detail-stats">
          <div class="detail-stat">
            <span class="detail-stat-label">Выживание</span>
            <span class="detail-stat-value">{{.Info.HoursSurvived}}</span>
          </div>
          <div class="detail-stat">
            <span class="detail-stat-label">Состояние</span>
            <span class="detail-stat-value">
              {{.Info.Condition}}%
              <span class="condition-bar {{.CondClass}}"><span class="condition-fill" style="width:{{.CondWidth}}"></span></span>
            </span>
          </div>
          <div class="detail-stat">
            <span class="detail-stat-label">Мир исследован</span>
            <span class="detail-stat-value">{{.Info.WorldExplored}}%</span>
          </div>
          <div class="detail-stat">
            <span class="detail-stat-label">Персонаж</span>
            <span class="detail-stat-value">{{.Info.Persona}}</span>
          </div>
        </div>

        <div class="detail-meta">
          Бэкап: {{.BackupDateFmt}}, {{.BackupTimeFmt}}<br>
          Сохранение: {{.Info.DisplayName}} &middot; {{.Info.Timestamp}}
        </div>

        <button class="btn-restore" onclick="showModal()">&#8635; ВОССТАНОВИТЬ</button>
      </div>
    </div>
  </div>
</div>

<div id="modal" class="modal-overlay hidden">
  <div class="modal-box">
    <h2>ВОССТАНОВИТЬ СОХРАНЕНИЕ?</h2>
    <p>{{.Info.Region}} &middot; {{.Info.HoursSurvived}}</p>
    <p>Бэкап от {{.BackupTimeFmt}}</p>
    <p class="modal-warning">Текущее сохранение будет перезаписано.</p>
    <div class="modal-buttons">
      <button class="btn-cancel" onclick="hideModal()">ОТМЕНА</button>
      <button class="btn-restore" onclick="doRestore()">ВОССТАНОВИТЬ</button>
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
		BackupDateFmt: formatDateRu(bDate),
		BackupTimeFmt: bTime,
		CondClass:     condClass(info.Condition),
		CondWidth:     condWidth(info.Condition),
		ID:            id,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	detailTmpl.Execute(w, data)
}
