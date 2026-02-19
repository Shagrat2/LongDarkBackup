//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd
// +build darwin dragonfly freebsd linux netbsd openbsd

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
