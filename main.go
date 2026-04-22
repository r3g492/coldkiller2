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
	"coldkiller2/util"
	"time"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var lastLog = time.Now()
var showSplash = true
var splashTimer float32 = 0.0
var showInitMenu = true
var intermission = false
var intermissionTimer float32 = 0.0
var intermissionLoadStep = 0
var paused = false
var stageWonDelay float32 = -1

func main() {
	// setting
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowUndecorated)
	rl.InitWindow(0, 0, "coldkiller2")
	defer rl.CloseWindow()

	rl.SetTargetFPS(144)
	rl.SetExitKey(0)

	w, h := VirtualWidth, VirtualHeight
	rl.DisableCursor()

	rl.InitAudioDevice()
	sound.Init()
	model.Init()
	stage.InitStages()

	initRender()
	defer unloadRender()

	splashTex := util.LoadTextureFromEmbedded("raylib_144x144.png")
	defer rl.UnloadTexture(splashTex)

	floorTex := generateFloorTile()
	const floorTiles = 128
	floorMesh := rl.GenMeshPlane(400, 400, floorTiles, floorTiles)
	tcCount := int(floorMesh.VertexCount) * 2
	tc := unsafe.Slice(floorMesh.Texcoords, tcCount)
	for i := range tc {
		tc[i] *= floorTiles
	}
	rl.UpdateMeshBuffer(floorMesh, 1, unsafe.Slice((*byte)(unsafe.Pointer(floorMesh.Texcoords)), tcCount*4), 0)
	floorModel := rl.LoadModelFromMesh(floorMesh)
	rl.SetMaterialTexture(floorModel.Materials, rl.MapDiffuse, floorTex)
	defer rl.UnloadTexture(floorTex)
	defer rl.UnloadModel(floorModel)

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
		mouseLocation := virtualMousePosition()
		ip := input.ReadInput(keyMap)
		ip.MouseLocation = mouseLocation
		log(mouseLocation, dt, player)

		if showSplash {
			showSplash = !doSplash(&splashTimer, splashTex, dt, w, h)
			continue
		}

		if ip.EndGamePressed {
			if showInitMenu {
				break
			} else if intermission {
				intermission = false
				intermissionTimer = 0
				intermissionLoadStep = 0
				intermissionUpgrades[0] = upgradeNone
				intermissionUpgrades[1] = upgradeNone
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
			if doInitMenu(stageManager, w, h) {
				break
			}
			drawInputOverlay(w, h, ip, keyMap, true)
			continue
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
				keyMap,
			)
			continue
		}

		if stageManager.StageLost() {
			stageManager.Difficulty = 1
			rl.StopSound(sound.Track)
			rl.PlaySound(sound.YouLose)
			player.ResetStats()
			showInitMenu = true
		}

		if !rl.IsSoundPlaying(sound.Track) {
			rl.PlaySound(sound.Track)
		}

		if paused {
			action := doPauseMenu(w, h, ip, keyMap)
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

		// enemy
		stageManager.Mutate(dt)
		var ebc = enemyManager.Mutate(dt, player, structureManager)
		enemyManager.ProcessAnimation(dt, player)

		if stageManager.StageWon() && player.IsAlive() {
			if stageWonDelay < 0 {
				stageWonDelay = 1.0
			}
			stageWonDelay -= dt
			if stageWonDelay <= 0 {
				stageWonDelay = -1
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
		} else {
			stageWonDelay = -1
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
		blastManager.AddBlasts(enemyManager.BlastBuffer)
		blastManager.Mutate(dt)

		sight.UpdateSight(
			blastManager,
			bulletManager,
			enemyManager,
			structureManager,
			player,
		)

		beginFrame()
		rl.ClearBackground(rl.Gray)

		rl.BeginMode3D(player.Camera)
		rl.DrawModel(floorModel, rl.Vector3{Y: -2}, 1.0, rl.White)
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
		drawPlayerStats(player)
		enemyManager.DrawUi(player)
		drawCursor(mouseLocation, player)

		endFrame()
	}
}
