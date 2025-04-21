package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
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

		fFrom := filepath.Join(fBackupIDFolder, fFileName)
		fTo := filepath.Join(DataFolder, fFileName)

		return Copy(fFrom, fTo)
	})
	if err != nil {
		http.Error(w, "Failed to restore: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
