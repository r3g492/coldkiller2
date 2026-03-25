package main

import (
	"coldkiller2/blast"
	"coldkiller2/bullet"
	"coldkiller2/enemy"
	"coldkiller2/input"
	"coldkiller2/killer"
	"coldkiller2/sight"
	"coldkiller2/sound"
	"coldkiller2/stage"
	"coldkiller2/structure"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var lastLog = time.Now()
var lastScore int

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowUndecorated)
	rl.InitWindow(0, 0, "coldkiller2")

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
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	sound.Init()

	rl.SetTargetFPS(144)
	keyMap := input.DefaultWASD()

	bulletManager := bullet.CreateManager()
	blastManager := blast.CreateManager()
	structureManager := structure.CreateManager(structure.RADIUS)
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

	initNewGame(
		bulletManager,
		blastManager,
		structureManager,
		player,
		enemyManager,
		stageManager,
	)

	rl.DisableCursor()

	btnWidth, btnHeight := float32(400), float32(80)
	btnRect := rl.Rectangle{
		X:      float32(w)/2 - btnWidth/2,
		Y:      float32(h)/2 + 60,
		Width:  btnWidth,
		Height: btnHeight,
	}
	difficulty := 0
	sqSize := float32(50)
	spacing := float32(20)
	diffY := float32(h)/2 - 80

	for !rl.WindowShouldClose() {
		if rl.IsCursorHidden() {
			rl.EnableCursor()
		}
		ip := input.ReadInput(keyMap)
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		diffText := fmt.Sprintf("Difficulty: %d", difficulty)
		fontSize := int32(40)
		textWidth := rl.MeasureText(diffText, fontSize)
		textX := float32(w)/2 - float32(textWidth)/2

		rl.DrawText(diffText, int32(textX), int32(diffY+10), fontSize, rl.RayWhite)

		minusRect := rl.Rectangle{X: textX - sqSize - spacing, Y: diffY, Width: sqSize, Height: sqSize}
		plusRect := rl.Rectangle{X: textX + float32(textWidth) + spacing, Y: diffY, Width: sqSize, Height: sqSize}

		if drawButton(minusRect, "-", rl.DarkGray, rl.Gray, rl.White) && difficulty > 0 {
			difficulty--
		}
		if drawButton(plusRect, "+", rl.DarkGray, rl.Gray, rl.White) {
			difficulty++
		}

		buttonText := "Start Game"
		if drawButton(btnRect, buttonText, rl.Red, rl.Red, rl.Red) {
			stageManager.SetDifficulty(difficulty)
			stageManager.ResetScore()
			break
		}
		drawInputOverlay(w, h, ip, keyMap)
		rl.EndDrawing()
		continue
	}

	showLostMenu := false
	lost := false
	intermission := false
	var intermissionTimer float32 = 0.0

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		mouseLocation := rl.GetMousePosition()
		log(mouseLocation, dt, player)
		ip := input.ReadInput(keyMap)

		if showLostMenu {
			if rl.IsCursorHidden() {
				rl.EnableCursor()
			}

			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)

			diffText := fmt.Sprintf("Difficulty: %d", difficulty)
			fontSizeDiff := int32(40)
			textWidthDiff := rl.MeasureText(diffText, fontSizeDiff)
			textX := float32(w)/2 - float32(textWidthDiff)/2

			rl.DrawText(diffText, int32(textX), int32(diffY+10), fontSizeDiff, rl.RayWhite)

			minusRect := rl.Rectangle{X: textX - sqSize - spacing, Y: diffY, Width: sqSize, Height: sqSize}
			plusRect := rl.Rectangle{X: textX + float32(textWidthDiff) + spacing, Y: diffY, Width: sqSize, Height: sqSize}

			if drawButton(minusRect, "-", rl.DarkGray, rl.Gray, rl.White) && difficulty > 0 {
				difficulty--
			}
			if drawButton(plusRect, "+", rl.DarkGray, rl.Gray, rl.White) {
				difficulty++
			}

			buttonText := "Restart Game"
			xpText := fmt.Sprintf("Score: %d", stageManager.StageWon)

			fontSizeScore := int32(60)
			textWidthScore := rl.MeasureText(xpText, fontSizeScore)
			rl.DrawText(xpText, int32(w)/2-textWidthScore/2, int32(h)/2-180, fontSizeScore, rl.Red)

			if drawButton(btnRect, buttonText, rl.Red, rl.Red, rl.Red) || ip.ResetGamePressed {
				showLostMenu = false

				stageManager.SetDifficulty(difficulty)
				stageManager.ResetScore()

				if lost {
					initNewGame(
						bulletManager,
						blastManager,
						structureManager,
						player,
						enemyManager,
						stageManager,
					)
					lost = false
				}
			}

			rl.EndDrawing()
			continue
		}

		if gameLost(player) {
			rl.StopSound(sound.Track)
			rl.PlaySound(sound.YouLose)
			showLostMenu = true
			lost = true
			lastScore = bulletManager.PlayerXp
			initNewGame(
				bulletManager,
				blastManager,
				structureManager,
				player,
				enemyManager,
				stageManager,
			)
			difficulty = stageManager.Difficulty
		}

		if !rl.IsSoundPlaying(sound.Track) && !showLostMenu {
			rl.PlaySound(sound.Track)
		}

		if !rl.IsCursorHidden() {
			rl.DisableCursor()
		}

		if ip.EndGamePressed {
			showLostMenu = true
		}

		if ip.ResetGamePressed {
			initNewGame(
				bulletManager,
				blastManager,
				structureManager,
				player,
				enemyManager,
				stageManager,
			)
		}

		if intermission {
			intermissionTimer += dt

			rl.BeginDrawing()
			rl.ClearBackground(rl.DarkGray)

			stageText := fmt.Sprintf("Score: %d", stageManager.StageWon)
			stageSize := int32(40)
			stageWidth := rl.MeasureText(stageText, stageSize)
			rl.DrawText(stageText, int32(w)/2-stageWidth/2, int32(h)/2-80, stageSize, rl.RayWhite)

			diffInfoText := fmt.Sprintf("Current Difficulty: %d", stageManager.Difficulty)
			diffInfoSize := int32(30)
			diffInfoWidth := rl.MeasureText(diffInfoText, diffInfoSize)
			rl.DrawText(diffInfoText, int32(w)/2-diffInfoWidth/2, int32(h)/2-30, diffInfoSize, rl.LightGray)

			timeLeft := 1.0 - intermissionTimer
			if timeLeft < 0 {
				timeLeft = 0
			}
			timerText := fmt.Sprintf("%.1f...", timeLeft)
			timerSize := int32(30)
			timerWidth := rl.MeasureText(timerText, timerSize)
			rl.DrawText(timerText, int32(w)/2-timerWidth/2, int32(h)/2+40, timerSize, rl.Gray)

			if intermissionTimer >= 1.0 {
				intermission = false
				intermissionTimer = 0
				stageManager.GenerateNewStage()
				initNewGame(
					bulletManager,
					blastManager,
					structureManager,
					player,
					enemyManager,
					stageManager,
				)
			}

			rl.EndDrawing()
			continue
		}

		// enemy
		var ebc = enemyManager.Mutate(dt, player, structureManager)
		enemyManager.ProcessAnimation(dt, player)

		if gameWon(enemyManager) {
			intermission = true
			rl.PlaySound(sound.ThreeTwoOne)
			initNewGame(
				bulletManager,
				blastManager,
				structureManager,
				player,
				enemyManager,
				stageManager,
			)
			stageManager.ScoreUp()
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
		// rl.NewColor(10, 10, 15, 255)
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
		enemyManager.DrawUi(player)
		drawInputOverlay(w, h, ip, keyMap)
		drawCursor(mouseLocation, player)

		rl.EndDrawing()
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
	)
	stageManager.CreateNewStage(player.Position)
}

func log(mouseLocation rl.Vector2, dt float32, player *killer.Killer) {
	if time.Since(lastLog) >= 1000*time.Millisecond {
		msg := fmt.Sprintf("mouseLocation=%v, dt=%v", mouseLocation, dt)
		fmt.Println(msg)
		fmt.Println(player.MoveDirection)
		lastLog = time.Now()
	}
}

func drawCursor(mouseLocation rl.Vector2, player *killer.Killer) {
	mouseRay := rl.GetScreenToWorldRay(mouseLocation, player.Camera)

	t := float32(0.0)
	if mouseRay.Direction.Y != 0 {
		t = (player.Position.Y - mouseRay.Position.Y) / mouseRay.Direction.Y
	}

	target3D := rl.Vector3{
		X: mouseRay.Position.X + mouseRay.Direction.X*t,
		Y: player.Position.Y,
		Z: mouseRay.Position.Z + mouseRay.Direction.Z*t,
	}

	rl.BeginMode3D(player.Camera)

	rl.DrawLine3D(player.Position, target3D, rl.Green)

	rl.DrawSphere(target3D, 0.1, rl.Green)

	rl.EndMode3D()
}

func gameLost(player *killer.Killer) bool {
	return !player.IsAlive() && player.ActionTimeLeft <= 0
}

func gameWon(enemyManager *enemy.Manager) bool {
	return enemyManager.AliveEnemyCount == 0
}

func drawButton(rect rl.Rectangle, text string, baseColor, hoverColor, textColor rl.Color) bool {
	mousePoint := rl.GetMousePosition()
	isHovered := rl.CheckCollisionPointRec(mousePoint, rect)
	isPressed := isHovered && (rl.IsMouseButtonDown(rl.MouseLeftButton) || rl.IsMouseButtonReleased(rl.MouseLeftButton))

	currentColor := baseColor
	if isPressed {
		currentColor = rl.Maroon
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			return true
		}
	} else if isHovered {
		currentColor = hoverColor
	}

	rl.DrawRectangleRec(rect, rl.Fade(currentColor, 0.3))
	rl.DrawRectangleLinesEx(rect, 3, currentColor)

	fontSize := int32(30)
	textWidth := rl.MeasureText(text, fontSize)
	textX := int32(rect.X + (rect.Width / 2) - float32(textWidth/2))
	textY := int32(rect.Y + (rect.Height / 2) - float32(fontSize/2))

	rl.DrawText(text, textX, textY, fontSize, textColor)

	return false
}

func drawInputOverlay(w, h int, ip input.Input, keyMap input.KeyMap) {
	const (
		keySize      = 45
		spacing      = 6
		fontSize     = 18
		descSize     = 9
		rightMargin  = 40
		bottomMargin = 90
	)

	labelUp := input.GetKeyName(keyMap.Up)
	labelLeft := input.GetKeyName(keyMap.Left)
	labelDown := input.GetKeyName(keyMap.Down)
	labelRight := input.GetKeyName(keyMap.Right)
	labelReload := input.GetKeyName(keyMap.Reload)
	labelReset := input.GetKeyName(keyMap.ResetGame)
	labelEnd := input.GetKeyName(keyMap.EndGame)

	totalWidth := (keySize * 7) + (spacing * 6) + 20
	baseX := float32(w) - float32(totalWidth) - rightMargin
	baseY := float32(h) - (keySize * 2) - spacing - bottomMargin

	drawKey := func(x, y float32, width float32, label string, desc string, active bool) {
		rect := rl.Rectangle{X: x, Y: y, Width: width, Height: keySize}

		alpha := uint8(100)
		if active {
			alpha = 180
		}

		bgColor := rl.NewColor(30, 30, 30, alpha)
		borderCol := rl.Fade(rl.Gray, 0.4)
		textCol := rl.Fade(rl.LightGray, 0.9)
		descCol := rl.Fade(rl.Gray, 0.7)

		if active {
			bgColor = rl.NewColor(230, 41, 55, 160)
			borderCol = rl.Red
			textCol = rl.White
			descCol = rl.Fade(rl.White, 0.8)
		}

		rl.DrawRectangleRec(rect, bgColor)
		rl.DrawRectangleLinesEx(rect, 1, borderCol)

		tw := rl.MeasureText(label, fontSize)
		rl.DrawText(label, int32(x+width/2)-tw/2, int32(y+keySize/2)-fontSize/2-4, fontSize, textCol)

		dtw := rl.MeasureText(desc, descSize)
		rl.DrawText(desc, int32(x+width/2)-dtw/2, int32(y+keySize)-descSize-6, descSize, descCol)
	}

	drawKey(baseX, baseY, keySize, labelEnd, "END", rl.IsKeyDown(keyMap.EndGame))
	drawKey(baseX+keySize+spacing, baseY, keySize, labelReset, "RESET", ip.ResetGamePressed)
	drawKey(baseX+(keySize+spacing)*3, baseY, keySize, labelUp, "UP", ip.MoveUp)

	currX := baseX + (keySize+spacing)*2
	drawKey(currX, baseY+keySize+spacing, keySize, labelLeft, "LEFT", ip.MoveLeft)
	currX += keySize + spacing

	drawKey(currX, baseY+keySize+spacing, keySize, labelDown, "DOWN", ip.MoveDown)
	currX += keySize + spacing

	drawKey(currX, baseY+keySize+spacing, keySize, labelRight, "RIGHT", ip.MoveRight)
	currX += keySize + spacing

	drawKey(currX, baseY+keySize+spacing, keySize, labelReload, "RELOAD", ip.ReloadPressed)
	currX += keySize + spacing

	drawKey(currX, baseY+keySize+spacing, keySize+30, "LMB", "SHOOT", rl.IsMouseButtonDown(keyMap.PunchHold))
}
