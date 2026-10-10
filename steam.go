//go:build steam

package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/go-steamworks"
)

// Built only with -tags steam. The Steam runtime library must be loadable:
// on Windows steam_api64.dll has to sit next to the exe, on Linux/macOS it is
// embedded by go-steamworks.

const steamAppID = 4690070

// achievementGameCleared must match the API Name set in Steamworks.
const achievementGameCleared = "ACH_GAME_CLEARED"

var steamReady bool

func steamInit() {
	if steamworks.RestartAppIfNecessary(steamAppID) {
		os.Exit(0)
	}
	if err := steamworks.Init(); err != nil {
		fmt.Println("steam: init failed, achievements disabled:", err)
		return
	}
	steamReady = true
}

func steamUpdate() {
	if steamReady {
		steamworks.RunCallbacks()
	}
}

func steamUnlock(name string) {
	if !steamReady {
		return
	}
	stats := steamworks.SteamUserStats()
	if achieved, ok := stats.GetAchievement(name); ok && achieved {
		return
	}
	if stats.SetAchievement(name) {
		stats.StoreStats()
	}
}
