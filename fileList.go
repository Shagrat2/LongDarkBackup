package main

import (
	"io/fs"
	"log"
	"path/filepath"
)

type TFileInfo struct {
	Path    string
	ModTime int64
}

type TFileList []TFileInfo

func NewListScan(Dir string) (ret TFileList, err error) {

	ret = make(TFileList, 0, 100)

	err = filepath.Walk(Dir, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			//log.Println("++++ " + path)
			return nil
		}

		ret = append(ret, TFileInfo{
			Path:    path,
			ModTime: info.ModTime().UnixMicro(),
		})

		return nil
	})

	return
}

func (t *TFileList) Index(fInfo TFileInfo) int {

	if t == nil {
		return -1
	}

	for ind, itm := range *t {

		if itm.Path != fInfo.Path {
			continue
		}

		if itm.ModTime == fInfo.ModTime {
			return ind
		}
	}

	return -1
}

func (t *TFileList) GetChanged(lst TFileList) TFileList {

	if t == nil {
		return nil
	}

	fFoundFiles := make(TFileList, 0, 100)

	for _, itm1 := range *t {

		if lst.Index(itm1) != -1 {
			continue
		}

		// Not found or change
		fFoundFiles = append(fFoundFiles, itm1)
	}

	return fFoundFiles
}

func (t *TFileList) Equal(lst TFileList) bool {

	if t == nil && lst == nil {
		return true
	}

	if t == nil || lst == nil {
		return false
	}

	// List 1
	cnt := 0
	for _, itm := range *t {
		if lst.Index(itm) != -1 {
			cnt++
		}
	}
	if cnt != len(*t) {
		return false
	}

	// List 2
	cnt = 0
	for _, itm := range lst {
		if t.Index(itm) != -1 {
			cnt++
		}
	}

	return cnt == len(lst)
}
