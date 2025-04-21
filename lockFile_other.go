//go:build !wasm
// +build !wasm

package main

import (
	"os"
	"syscall"
)

func LockFile(in *os.File) error {
	return syscall.Flock(int(in.Fd()), syscall.LOCK_EX)
}

func UnLockFile(in *os.File) error {
	return syscall.Flock(int(in.Fd()), syscall.LOCK_UN)
}
