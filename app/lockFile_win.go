//go:build windows
// +build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"
)

type lockType uint32

const (
	readLock  lockType = 0
	writeLock lockType = windows.LOCKFILE_EXCLUSIVE_LOCK
)

const (
	reserved = 0
	allBytes = ^uint32(0)
)

func LockFile(in *os.File) error {
	return windows.LockFileEx(windows.Handle(in.Fd()), uint32(writeLock), reserved, allBytes, allBytes, nil)
}

func UnLockFile(in *os.File) error {
	return windows.UnlockFileEx(windows.Handle(in.Fd()), reserved, allBytes, allBytes, nil)
}
