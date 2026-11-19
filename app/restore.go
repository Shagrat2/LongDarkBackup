package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type restoreResult struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// Copy replaces dstFile atomically so a failed restore never corrupts the destination.
func Copy(srcFile, dstFile string) error {
	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.CreateTemp(filepath.Dir(dstFile), filepath.Base(dstFile)+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := out.Name()
	defer os.Remove(tmpName)

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Sync(); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, dstFile)
}

func DoRestoreHandler(w http.ResponseWriter, r *http.Request) {

	// Restore changes save files, so it must be an explicit same-origin POST
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if !isSameOrigin(r) {
		writeRestoreResult(w, http.StatusForbidden, restoreResult{Error: "forbidden"})
		return
	}

	ID := filepath.Base(r.URL.Path)

	fBackupIDFolder := filepath.Join(GetBackupFolder(), ID)

	err := filepath.WalkDir(fBackupIDFolder, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip dir
		if d.IsDir() {
			return nil
		}

		fFileName := filepath.Base(path)

		// Skip cache files
		if strings.HasPrefix(fFileName, "info.") {
			return nil
		}

		fFrom := filepath.Join(fBackupIDFolder, fFileName)
		fTo := filepath.Join(DataFolder, fFileName)

		return Copy(fFrom, fTo)
	})

	if err != nil {
		writeRestoreResult(w, http.StatusOK, restoreResult{Error: err.Error()})
		return
	}

	writeRestoreResult(w, http.StatusOK, restoreResult{OK: true})
}

func writeRestoreResult(w http.ResponseWriter, status int, res restoreResult) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("Restore response error:", err)
	}
}
