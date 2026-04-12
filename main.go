package main

import (
	"coldkiller2/blast"
	"coldkiller2/bullet"
	"coldkiller2/enemy"
	"coldkiller2/input"
	"coldkiller2/killer"
	"coldkiller2/model"
	"coldkiller2/sight"
	"coldkiller2/sound"
	"coldkiller2/stage"
	"coldkiller2/structure"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var lastLog = time.Now()
var showInitMenu = true
var intermission = false
var intermissionTimer float32 = 0.0
var paused = false

func main() {
	// setting
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowUndecorated)
	rl.InitWindow(0, 0, "coldkiller2")
	defer rl.CloseWindow()

	rl.SetTargetFPS(144)
	rl.SetExitKey(0)

	// TODO: monitor change feature?
	w, h := setMonitor()
	rl.DisableCursor()

	rl.InitAudioDevice()
	sound.Init()
	model.Init()
	stage.InitStages()

	keyMap := input.DefaultWASD()

	bulletManager := bullet.CreateManager()
	blastManager := blast.CreateManager()
	structureManager := structure.CreateManager()
	enemyManager := enemy.CreateManager()
	stageManager := stage.CreateManager()
	stageManager.HighestBeaten = stage.LoadProgress()
	player := killer.Create()
	defer unloadGame(
		bulletManager,
		blastManager,
		structureManager,
		player,
		enemyManager,
		stageManager,
	)

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		mouseLocation := rl.GetMousePosition()
		ip := input.ReadInput(keyMap)
		log(mouseLocation, dt, player)

		if ip.EndGamePressed {
			if showInitMenu {
				break
			} else if intermission {
				intermission = false
				intermissionTimer = 0
				showInitMenu = true
				rl.StopSound(sound.Track)
			} else {
				paused = !paused
			}
		}

		if stageManager.GameWon() {
			if doGameWon(w, h) {
				break
			}
			continue
		}

		if showInitMenu {
			if doInitMenu(
				stageManager,
				bulletManager,
				blastManager,
				structureManager,
				player,
				enemyManager,
				w,
				h,
			) {
				break
			}
			drawInputOverlay(w, h, ip, keyMap)
			continue
		}

		if stageManager.StageLost() {
			rl.StopSound(sound.Track)
			rl.PlaySound(sound.YouLose)
			showInitMenu = true
		}

		if !rl.IsSoundPlaying(sound.Track) {
			rl.PlaySound(sound.Track)
		}

		if paused {
			if rl.IsCursorHidden() {
				rl.EnableCursor()
			}
			action := doPauseMenu(w, h)
			if action == pauseResume {
				paused = false
			} else if action == pauseQuitToMenu {
				paused = false
				showInitMenu = true
				rl.StopSound(sound.Track)
			}
			continue
		}

		if !rl.IsCursorHidden() {
			rl.DisableCursor()
		}

		if intermission {
			doIntermission(
				dt,
				stageManager,
				bulletManager,
				blastManager,
				structureManager,
				player,
				enemyManager,
				w,
				h,
			)
			continue
		}

		// enemy
		var ebc = enemyManager.Mutate(dt, player, structureManager)
		enemyManager.ProcessAnimation(dt, player)

		if stageManager.StageWon() {
			intermission = true
			if stageManager.Difficulty > stageManager.HighestBeaten {
				stageManager.HighestBeaten = stageManager.Difficulty
				stage.SaveProgress(stageManager.HighestBeaten)
			}
			stageManager.Difficulty++
			if !stageManager.GameWon() {
				// rl.PlaySound(sound.ThreeTwoOne)
			}
		}

		// player
		bc := player.Mutate(ip, dt, enemyManager.GetBoundingBoxes(), structureManager)
		player.ResolveAnimation()
		player.PlanAnimate(dt)
		player.Animate()

		// bullet
		bulletManager.KillerBulletCreate(bc)
		bulletManager.EnemyBulletCreate(ebc)
		bulletBlasts := bulletManager.Mutate(dt, player, enemyManager.Enemies, structureManager)

		blastManager.AddBlasts(bulletBlasts)
		blastManager.Mutate(dt)

		sight.UpdateSight(
			blastManager,
			bulletManager,
			enemyManager,
			structureManager,
			player,
		)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Gray)

		rl.BeginMode3D(player.Camera)
		sight.DrawSolidShadows(player.Position, structureManager)
		player.Draw3D()
		enemyManager.Draw3D(player)
		bulletManager.Draw3D()
		structureManager.Draw3D(player.Position)
		blastManager.Draw3D()
		rl.EndMode3D()

		player.DrawUi()
		player.DrawHitFlash()
		drawEnemyCount(w, enemyManager.AliveCount())
		enemyManager.DrawUi(player)
		drawInputOverlay(w, h, ip, keyMap)
		drawCursor(mouseLocation, player)

		rl.EndDrawing()
	}
}
