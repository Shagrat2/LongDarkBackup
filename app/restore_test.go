package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func restoreFixture(t *testing.T) (id string, backups string, data string) {
	t.Helper()
	backups = t.TempDir()
	data = t.TempDir()
	SetBackupFolder(backups)
	DataFolder = data

	id = "2026-01-01-00-00-00"
	if err := os.MkdirAll(filepath.Join(backups, id), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(backups, id, "sandbox"), []byte("save"), 0644); err != nil {
		t.Fatal(err)
	}
	return id, backups, data
}

func TestCopyKeepsDestinationOnFailure(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "sandbox")
	if err := os.WriteFile(dst, []byte("current save"), 0644); err != nil {
		t.Fatal(err)
	}

	src := filepath.Join(t.TempDir(), "missing")
	if err := Copy(src, dst); err == nil {
		t.Fatal("expected copy error")
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "current save" {
		t.Fatalf("destination was modified: %q", got)
	}
}

func TestRestoreRequiresPost(t *testing.T) {
	id, _, data := restoreFixture(t)

	req := httptest.NewRequest(http.MethodGet, "/restore/"+id, nil)
	w := httptest.NewRecorder()
	DoRestoreHandler(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}
	if _, err := os.Stat(filepath.Join(data, "sandbox")); !os.IsNotExist(err) {
		t.Fatal("GET request must not restore files")
	}
}

func TestRestoreRejectsCrossOrigin(t *testing.T) {
	id, _, data := restoreFixture(t)

	req := httptest.NewRequest(http.MethodPost, "/restore/"+id, nil)
	req.Header.Set("Origin", "http://evil.example")
	w := httptest.NewRecorder()
	DoRestoreHandler(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
	if _, err := os.Stat(filepath.Join(data, "sandbox")); !os.IsNotExist(err) {
		t.Fatal("cross-origin request must not restore files")
	}
}

func TestRestoreSameOrigin(t *testing.T) {
	id, _, data := restoreFixture(t)

	req := httptest.NewRequest(http.MethodPost, "/restore/"+id, nil)
	w := httptest.NewRecorder()
	DoRestoreHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	got, err := os.ReadFile(filepath.Join(data, "sandbox"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "save" {
		t.Fatalf("restored content = %q", got)
	}
}
