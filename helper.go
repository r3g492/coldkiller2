package main

import (
	"coldkiller2/blast"
	"coldkiller2/bullet"
	"coldkiller2/enemy"
	"coldkiller2/killer"
	"coldkiller2/stage"
	"coldkiller2/structure"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func setMonitor() (int, int) {
	targetMonitor := 0
	monitorCount := rl.GetMonitorCount()
	if targetMonitor >= monitorCount {
		targetMonitor = 0
	}
	w := rl.GetMonitorWidth(targetMonitor)
	h := rl.GetMonitorHeight(targetMonitor)
	monitorPos := rl.GetMonitorPosition(targetMonitor)
	rl.SetWindowPosition(int(monitorPos.X), int(monitorPos.Y))
	rl.SetWindowSize(w, h)
	return w, h
}

func log(mouseLocation rl.Vector2, dt float32, player *killer.Killer) {
	if time.Since(lastLog) >= 1000*time.Millisecond {
		msg := fmt.Sprintf("mouseLocation=%v, dt=%v", mouseLocation, dt)
		fmt.Println(msg)
		fmt.Println(player.MoveDirection)
		lastLog = time.Now()
	}
}

func unloadGame(
	bulletManager *bullet.Manager,
	blastManager *blast.Manager,
	structureManager *structure.Manager,
	player *killer.Killer,
	enemyManager *enemy.Manager,
	stageManager *stage.Manager,
) {
	bulletManager.Unload()
	blastManager.Unload()
	structureManager.Unload()
	player.Unload()
	enemyManager.Unload()
	stageManager.Unload()
}

func initNewGame(
	bulletManager *bullet.Manager,
	blastManager *blast.Manager,
	structureManager *structure.Manager,
	player *killer.Killer,
	enemyManager *enemy.Manager,
	stageManager *stage.Manager,
) {
	bulletManager.Init()
	blastManager.Init()
	structureManager.Init()
	player.Init()
	enemyManager.Init(player)
	stageManager.Init(
		structureManager,
		enemyManager,
		player,
	)
}
