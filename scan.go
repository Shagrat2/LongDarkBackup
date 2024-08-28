package main

import (
	// "crypto/md5"
	// "fmt"
	// "io"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

var (
	fDataDir   = "/Users/ivan/.local/share/Hinterland"
	fBackupDir = "/Users/ivan/Documents/LongDarkBackup"
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
	err = syscall.Flock(int(in.Fd()), syscall.LOCK_EX)
	if err != nil {
		log.Println("Error lock src", err)
		return
	}
	defer syscall.Flock(int(in.Fd()), syscall.LOCK_UN)

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
	itms, err := NewListScan(fDataDir)
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
	try := 0
	for {
		if try >= 30 {
			log.Println("Error wait changed")
			return
		}

		time.Sleep(1 * time.Second)

		tmpList, err := NewListScan(fDataDir)
		if err != nil {
			log.Println(err)
			return
		}

		// get changed
		fChFiles := tmpList.GetChanged(LastTree)
		if fFoundFiles.Equal(fChFiles) {
			break
		}

		fFoundFiles = tmpList
		try++
	}

	//=========
	LastTree = itms

	// Backup
	fFolderName := time.Now().Format("2006-01-02-15-04-05")
	fToFolder := filepath.Join(fBackupDir, fFolderName)
	os.MkdirAll(fToFolder, 0755)

	log.Println("######")
	for _, itm := range fFoundFiles {
		destFile := filepath.Join(fToFolder, filepath.Base(itm.Path))
		copyFileContents(itm.Path, destFile)

		log.Println("### " + itm.Path)
	}

	//! Save file list
}

func DoGarbage() {
	//!
}
