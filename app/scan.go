package main

import (
	// "crypto/md5"
	// "fmt"
	// "io"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// func FileMD5(path string) string {
// 	h := md5.New()
// 	f, err := os.Open(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	_, err = io.Copy(h, f)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return fmt.Sprintf("%x", h.Sum(nil))
// }

var LastTree TFileList

func copyFileContents(src, dst string) (err error) {

	// open src
	in, err := os.OpenFile(src, os.O_RDONLY, 0644)
	if err != nil {
		log.Println("Error open src file: ", err)
		return
	}
	defer in.Close()

	// Lock src
	err = LockFile(in)
	if err != nil {
		log.Println("Error lock src", err)
		return
	}
	defer UnLockFile(in)

	// Create dest
	out, err := os.Create(dst)
	if err != nil {
		log.Println("Open dst", err)
		return
	}
	defer out.Close()

	// Copy data
	if _, err = io.Copy(out, in); err != nil {
		log.Println("Error copy", err)
		return
	}

	// Save cache
	err = out.Sync()
	return
}

func DoScan() {

	// Scan files
	itms, err := NewListScan(DataFolder)
	if err != nil {
		log.Println(err)
		return
	}

	if LastTree == nil {
		LastTree = itms
		return
	}

	// Compare list
	fFoundFiles := itms.GetChanged(LastTree)

	// Not found
	if len(fFoundFiles) == 0 {
		return
	}

	// Wait finish change
	var stableList TFileList
	try := 0
	for {
		if try >= 30 {
			log.Println("Error wait changed")
			return
		}

		time.Sleep(1 * time.Second)

		tmpList, err := NewListScan(DataFolder)
		if err != nil {
			log.Println(err)
			return
		}

		// get changed
		fChFiles := tmpList.GetChanged(LastTree)
		if fFoundFiles.Equal(fChFiles) {
			stableList = tmpList
			break
		}

		fFoundFiles = fChFiles
		try++
	}

	// Backup only changed files
	fFolderName := time.Now().Format("2006-01-02-15-04-05")
	fToFolder := filepath.Join(GetBackupFolder(), fFolderName)
	if err := os.MkdirAll(fToFolder, 0755); err != nil {
		log.Println("Cannot create backup folder:", err)
		return
	}

	log.Println("######")
	failed := make(map[string]bool)
	for _, itm := range fFoundFiles {
		destFile := filepath.Join(fToFolder, filepath.Base(itm.Path))
		if err := copyFileContents(itm.Path, destFile); err != nil {
			log.Println("Backup error:", err)
			os.Remove(destFile)
			failed[itm.Path] = true
			continue
		}

		log.Println("### " + itm.Path)
	}

	if len(failed) > 0 {
		// Keep the successful files and retry only the failed ones next scan
		processed := make(TFileList, 0, len(stableList))
		for _, itm := range stableList {
			if !failed[itm.Path] {
				processed = append(processed, itm)
			}
		}
		LastTree = processed

		if len(failed) == len(fFoundFiles) {
			os.Remove(fToFolder)
		}
		return
	}

	// Mark changes as processed only after a complete backup
	LastTree = stableList
}

func DoGarbage() {
	//!
}
