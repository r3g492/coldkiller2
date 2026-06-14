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

var showOptions = false
var showCleared = false
var winSoundPlayed = false

var keyMap input.KeyMap

var kpsKills int
var kpsTimeSec float32
var kpsPrevAlive int
var kpsLatest float32

var audioPitch float32 = 1.0

const hitMarkerDuration = 0.25

var hitMarkerTimer float32

const testMode = false

var testConfig = struct {
	SlowDuration float32
	AmmoCapacity int32
	MoveSpeed    float32
	Stage        int
}{
	SlowDuration: 5.0,
	AmmoCapacity: 30,
	MoveSpeed:    12.0,
	Stage:        1,
}

func main() {
	loadConfig()

	initFlags := uint32(rl.FlagWindowResizable)
	switch currentConfig.WindowMode {
	case WindowModeBorderless:
		initFlags |= uint32(rl.FlagWindowUndecorated) | uint32(rl.FlagBorderlessWindowedMode)
	case WindowModeWindowed:
		// WindowModeWindowed: no extra flags — window gets decorations
	}
	rl.SetConfigFlags(initFlags)
	rl.InitWindow(0, 0, "Kill Per Second")
	defer rl.CloseWindow()
	defer restoreSystemUI()

	monitor := rl.GetCurrentMonitor()
	mw := rl.GetMonitorWidth(monitor)
	mh := rl.GetMonitorHeight(monitor)
	switch currentConfig.WindowMode {
	case WindowModeBorderless:
		rl.SetWindowPosition(0, 0)
		rl.SetWindowSize(mw, mh)
		hideSystemUI()
	case WindowModeWindowed:
		rw, rh := currentConfig.ResWidth, currentConfig.ResHeight
		rl.SetWindowSize(rw, rh)
		rl.SetWindowPosition((mw-rw)/2, (mh-rh)/2)
	}

	rl.SetTargetFPS(60)
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

	winTex := util.LoadTextureFromEmbedded("win.png")
	rl.GenTextureMipmaps(&winTex)
	rl.SetTextureFilter(winTex, rl.FilterTrilinear)
	defer rl.UnloadTexture(winTex)

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

	keyMap = input.DefaultWASD()
	if currentConfig.KeyBindings != nil {
		keyMap = *currentConfig.KeyBindings
	}

	bulletManager := bullet.CreateManager()
	blastManager := blast.CreateManager()
	structureManager := structure.CreateManager()
	enemyManager := enemy.CreateManager()
	stageManager := stage.CreateManager()
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
			} else if showOptions {
				showOptions = false
				showInitMenu = true
			} else if showCleared {
				showCleared = false
				showInitMenu = true
			} else if editing {
				editing = false
				showInitMenu = true
			} else {
				paused = !paused
			}
		}

		if editing {
			updateEditor(dt, mouseLocation)
			drawEditor(mouseLocation)
			continue
		}

		if stageManager.GameWon() {
			if !currentConfig.GameCleared {
				currentConfig.GameCleared = true
				saveConfig()
			}
			if !winSoundPlayed {
				rl.StopSound(sound.Track)
				rl.PlaySound(sound.Win)
				winSoundPlayed = true
			}
			if doGameWon(w, h, "Main Menu", winTex) {
				stageManager.Difficulty = 0
				showInitMenu = true
			}
			continue
		}

		if showCleared {
			if doGameWon(w, h, "Back", winTex) {
				showCleared = false
				showInitMenu = true
			}
			continue
		}

		if showInitMenu {
			kpsKills = 0
			kpsTimeSec = 0
			kpsPrevAlive = 0
			kpsLatest = 0
			resetCombo()
			action := doInitMenu(w, h)
			if action == initMenuExit {
				break
			}
			if action == initMenuStart {
				rl.StopSound(sound.Win)
				winSoundPlayed = false
				player.ResetStats()
				stageManager.Difficulty = 1
				showInitMenu = false
				intermission = true
			} else if action == initMenuTestStart {
				rl.StopSound(sound.Win)
				winSoundPlayed = false
				player.ResetStats()
				player.SlowTimeDuration = testConfig.SlowDuration
				player.SlowTimeLeft = testConfig.SlowDuration
				player.AmmoCapacity = testConfig.AmmoCapacity
				player.MoveSpeed = testConfig.MoveSpeed
				stageManager.Difficulty = testConfig.Stage
				showInitMenu = false
				intermission = true
			} else if action == initMenuEditor {
				showInitMenu = false
				openEditor()
			} else if action == initMenuContinue {
				loadProgress(player)
				stageManager.Difficulty = currentConfig.StageLevel + 1
				showInitMenu = false
				intermission = true
			} else if action == initMenuOptions {
				showInitMenu = false
				showOptions = true
			} else if action == initMenuShowCleared {
				showInitMenu = false
				showCleared = true
				rl.PlaySound(sound.Win)
			}
			continue
		}

		if showOptions {
			if doOptionsMenu(w, h) {
				showOptions = false
				showInitMenu = true
			}
			continue
		}

		if intermission {
			if paused {
				action := doPauseMenu(w, h, ip, keyMap, true)
				if action == pauseResume {
					paused = false
				} else if action == pauseQuitToMenu {
					paused = false
					intermission = false
					intermissionTimer = 0
					intermissionLoadStep = 0
					showInitMenu = true
					rl.StopSound(sound.Track)
				}
				continue
			}
			gaveUp := doIntermission(
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
			if gaveUp {
				intermission = false
				intermissionTimer = 0
				intermissionLoadStep = 0
				showInitMenu = true
				rl.StopSound(sound.Track)
			}
			continue
		}

		if stageManager.StageLost() {
			stageManager.Difficulty = 0
			rl.StopSound(sound.Track)
			rl.PlaySound(sound.YouLose)
			player.ResetStats()
			showInitMenu = true
		}

		if !rl.IsSoundPlaying(sound.Track) {
			rl.PlaySound(sound.Track)
		}

		if paused {
			action := doPauseMenu(w, h, ip, keyMap, true)
			if action == pauseResume {
				paused = false
			} else if action == pauseQuitToMenu {
				paused = false
				stageManager.Difficulty = 0
				showInitMenu = true
				rl.StopSound(sound.Track)
			}
			continue
		}

		if !rl.IsCursorHidden() {
			rl.DisableCursor()
		}

		worldDt := dt
		targetPitch := float32(1.0)
		if player.SlowTimeActive {
			worldDt = dt * 0.1
			targetPitch = 0.4
		}
		audioPitch += (targetPitch - audioPitch) * 10 * dt
		sound.SetGlobalPitch(audioPitch)

		// enemy
		stageManager.Mutate(worldDt)
		var ebc = enemyManager.Mutate(worldDt, player, structureManager)
		enemyManager.ProcessAnimation(worldDt, player)

		if kpsPrevAlive > enemyManager.AliveEnemyCount {
			kpsKills += kpsPrevAlive - enemyManager.AliveEnemyCount
		}
		kpsPrevAlive = enemyManager.AliveEnemyCount
		updateCombo(enemyManager.PlayerBulletKillCount+enemyManager.ExplosionKillCount, dt)
		if !stageManager.StageWon() {
			kpsTimeSec += worldDt
		}

		if stageManager.StageWon() && player.IsAlive() {
			if stageWonDelay < 0 {
				stageWonDelay = 1.0
			}
			stageWonDelay -= dt
			if stageWonDelay <= 0 {
				stageWonDelay = -1
				intermission = true
				stageKps := float32(0)
				if kpsTimeSec > 0 {
					stageKps = float32(kpsKills) / kpsTimeSec
				}
				kpsLatest = stageKps
				if stageKps > currentConfig.BestKps {
					currentConfig.BestKps = stageKps
				}
				saveProgress(stageManager.Difficulty, player)
				stageManager.Difficulty++
				kpsKills = 0
				kpsTimeSec = 0
				kpsPrevAlive = 0
				resetCombo()
			}
		} else {
			stageWonDelay = -1
		}

		// player
		bc := player.Mutate(ip, worldDt, dt, enemyManager.GetBoundingBoxes(), structureManager)
		player.ResolveAnimation()
		player.PlanAnimate(dt)
		player.Animate()

		// bullet
		muzzleBlasts := bulletManager.KillerBulletCreate(bc)
		bulletManager.EnemyBulletCreate(ebc)
		bulletBlasts, playerHits := bulletManager.Mutate(worldDt, player, enemyManager.Enemies, structureManager)
		if playerHits > 0 {
			hitMarkerTimer = hitMarkerDuration
		}
		if hitMarkerTimer > 0 {
			hitMarkerTimer -= dt
		}

		blastManager.AddBlasts(muzzleBlasts)
		blastManager.AddBlasts(bulletBlasts)
		blastManager.AddBlasts(enemyManager.BlastBuffer)
		blastManager.Mutate(worldDt)

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
		player.DrawSlowTimeVignette()
		drawKPS(w, 50)
		drawCombo(w)
		enemyManager.DrawUi(player)
		enemyManager.DrawOffscreenIndicators(player)
		drawEscOverlay()
		drawCursor(mouseLocation, player)
		drawHitMarker(mouseLocation)
		rl.DrawFPS(int32(w)-90, 10)

		endFrame()
	}
}
