//go:build linux

package main

import rl "github.com/gen2brain/raylib-go/raylib"

const windowModeConfigurable = false

func hideSystemUI() {
	if !rl.IsWindowFullscreen() {
		rl.ToggleFullscreen()
	}
}

func restoreSystemUI() {
	if rl.IsWindowFullscreen() {
		rl.ToggleFullscreen()
	}
}
