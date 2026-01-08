//go:build windows

package main

import (
	"syscall"
)

func attachConsole() {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procAttachConsole := modkernel32.NewProc("AttachConsole")
	const ATTACH_PARENT_PROCESS = ^uintptr(0) // -1
	procAttachConsole.Call(ATTACH_PARENT_PROCESS)
}
