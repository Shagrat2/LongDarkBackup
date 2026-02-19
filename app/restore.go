package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func DoRestoreHandler(w http.ResponseWriter, r *http.Request) {

	ID := filepath.Base(r.URL.Path)

	fBackupIDFolder := filepath.Join(BackupFolder, ID)

	err := filepath.WalkDir(fBackupIDFolder, func(path string, d os.DirEntry, err error) error {

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

	// POST → JSON response
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			fmt.Fprintf(w, `{"ok":false,"error":"%s"}`, err.Error())
		} else {
			fmt.Fprint(w, `{"ok":true}`)
		}
		return
	}

	// GET → redirect (backward compatible)
	if err != nil {
		http.Error(w, "Failed to restore: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
