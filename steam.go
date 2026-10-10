//go:build steam

package main

import (
	"fmt"
	"os"
	"time"

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
	fmt.Println("steam: initialized")
}

// Unlocks are queued and retried from steamUpdate: right after Init the
// user's stats are still being fetched, and SetAchievement fails until then.
var (
	pendingUnlocks = map[string]bool{}
	nextUnlockTry  time.Time
)

func steamUpdate() {
	if !steamReady {
		return
	}
	steamworks.RunCallbacks()
	if len(pendingUnlocks) == 0 || time.Now().Before(nextUnlockTry) {
		return
	}
	nextUnlockTry = time.Now().Add(time.Second)
	stats := steamworks.SteamUserStats()
	for name := range pendingUnlocks {
		if achieved, ok := stats.GetAchievement(name); ok && achieved {
			delete(pendingUnlocks, name)
			continue
		}
		if stats.SetAchievement(name) && stats.StoreStats() {
			fmt.Println("steam: unlocked", name)
			delete(pendingUnlocks, name)
		}
	}
}

func steamUnlock(name string) {
	if steamReady {
		pendingUnlocks[name] = true
	}
}
