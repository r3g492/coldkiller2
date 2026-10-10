//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

// Go builds a console-subsystem exe by default, so launching the game from
// Explorer/Steam opens a terminal window that fills with raylib's trace logs.
// Release builds should use -ldflags "-H=windowsgui" so no console is created
// at all; this is a fallback for builds made without that flag.
func init() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleProcessList := kernel32.NewProc("GetConsoleProcessList")
	freeConsole := kernel32.NewProc("FreeConsole")

	// Only detach when the console exists solely for this process. If the game
	// was started from cmd/PowerShell, the shell shares the console and we keep
	// it so logs stay visible.
	var pids [2]uint32
	n, _, _ := getConsoleProcessList.Call(uintptr(unsafe.Pointer(&pids[0])), uintptr(len(pids)))
	if n == 1 {
		freeConsole.Call()
	}
}
