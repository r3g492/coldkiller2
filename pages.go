package main

import (
	"coldkiller2/blast"
	"coldkiller2/bullet"
	"coldkiller2/enemy"
	"coldkiller2/input"
	"coldkiller2/killer"
	"coldkiller2/sound"
	"coldkiller2/stage"
	"coldkiller2/structure"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	btnWidth  = float32(400)
	btnHeight = float32(80)
	sqSize    = float32(50)
	spacing   = float32(20)
)

var showDiffError = false

type pauseAction int

const (
	pauseNone       pauseAction = iota
	pauseResume     pauseAction = iota
	pauseQuitToMenu pauseAction = iota
)

func doPauseMenu(w, h int) pauseAction {
	rl.BeginDrawing()
	rl.ClearBackground(rl.NewColor(0, 0, 0, 0))
	rl.DrawRectangle(0, 0, int32(w), int32(h), rl.NewColor(0, 0, 0, 160))

	titleText := "PAUSED"
	titleSize := int32(60)
	titleWidth := rl.MeasureText(titleText, titleSize)
	rl.DrawText(titleText, int32(w)/2-titleWidth/2, int32(h)/2-160, titleSize, rl.RayWhite)

	resumeRect := rl.Rectangle{
		X:      float32(w)/2 - btnWidth/2,
		Y:      float32(h)/2 - 40,
		Width:  btnWidth,
		Height: btnHeight,
	}
	quitRect := rl.Rectangle{
		X:      float32(w)/2 - btnWidth/2,
		Y:      float32(h)/2 + 60,
		Width:  btnWidth,
		Height: btnHeight,
	}

	result := pauseNone
	if drawButton(resumeRect, "Resume", rl.DarkGray, rl.Gray, rl.White) {
		result = pauseResume
	}
	if drawButton(quitRect, "Quit to Menu", rl.DarkGray, rl.Gray, rl.White) {
		result = pauseQuitToMenu
	}

	rl.EndDrawing()
	return result
}

func doIntermission(
	dt float32,
	stageManager *stage.Manager,
	bulletManager *bullet.Manager,
	blastManager *blast.Manager,
	structureManager *structure.Manager,
	player *killer.Killer,
	enemyManager *enemy.Manager,
	w int,
	h int,
) {
	intermissionTimer += dt

	rl.BeginDrawing()
	rl.ClearBackground(rl.DarkGray)

	labelText := "NEXT STAGE"
	labelSize := int32(18)
	labelWidth := rl.MeasureText(labelText, labelSize)
	rl.DrawText(labelText, int32(w)/2-labelWidth/2, int32(h)/2-70, labelSize, rl.Gray)

	diffInfoText := fmt.Sprintf("%d / %d", stageManager.Difficulty, len(stage.Stages))
	diffInfoSize := int32(40)
	diffInfoWidth := rl.MeasureText(diffInfoText, diffInfoSize)
	rl.DrawText(diffInfoText, int32(w)/2-diffInfoWidth/2, int32(h)/2-40, diffInfoSize, rl.RayWhite)

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
		stageManager.CreateNewStage(player.Position)
	}

	rl.EndDrawing()
}

func doInitMenu(
	stageManager *stage.Manager,
	bulletManager *bullet.Manager,
	blastManager *blast.Manager,
	structureManager *structure.Manager,
	player *killer.Killer,
	enemyManager *enemy.Manager,
	w int,
	h int,
) bool {
	startingDiffUpperBound := len(stage.Stages)
	maxAllowed := stageManager.HighestBeaten + 1
	if maxAllowed < 1 {
		maxAllowed = 1
	}
	if maxAllowed > startingDiffUpperBound {
		maxAllowed = startingDiffUpperBound
	}

	diffY := float32(h)/2 - 80

	rl.StopSound(sound.Track)
	if stageManager.Difficulty > startingDiffUpperBound {
		stageManager.Difficulty = startingDiffUpperBound
	}

	if rl.IsCursorHidden() {
		rl.EnableCursor()
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	for k := int32(rl.KeyZero); k <= int32(rl.KeyNine); k++ {
		if rl.IsKeyPressed(k) {
			digit := int(k - int32(rl.KeyZero))
			next := stageManager.Difficulty*10 + digit
			if next >= 1 && next <= maxAllowed {
				stageManager.Difficulty = next
			} else if digit >= 1 && digit <= maxAllowed {
				stageManager.Difficulty = digit
			}
			showDiffError = false
		}
	}
	if rl.IsKeyPressed(rl.KeyBackspace) {
		stageManager.Difficulty /= 10 // single digit → 0 (blank)
		showDiffError = false
	}

	blank := stageManager.Difficulty == 0

	tryStart := rl.IsKeyPressed(rl.KeyEnter)
	fontSizeDiff := int32(40)
	var diffText string
	cursorVisible := int(rl.GetTime()*2)%2 == 0
	if blank {
		if cursorVisible {
			diffText = "_"
		} else {
			diffText = " "
		}
	} else {
		diffText = fmt.Sprintf("%d", stageManager.Difficulty)
	}
	textWidthDiff := rl.MeasureText(diffText, fontSizeDiff)
	textX := float32(w)/2 - float32(textWidthDiff)/2
	diffColor := rl.RayWhite
	if blank {
		diffColor = rl.Gray
	}
	rl.DrawText(diffText, int32(textX), int32(diffY+10), fontSizeDiff, diffColor)

	rangeText := fmt.Sprintf("type between 1 - %d", maxAllowed)
	if maxAllowed < startingDiffUpperBound {
		rangeText += fmt.Sprintf("  (play to unlock)")
	}
	rangeSize := int32(20)
	rangeWidth := rl.MeasureText(rangeText, rangeSize)
	rl.DrawText(rangeText, int32(w)/2-rangeWidth/2, int32(diffY)-30, rangeSize, rl.DarkGray)

	startRect := rl.Rectangle{
		X:      float32(w)/2 - btnWidth/2,
		Y:      float32(h)/2 + 60,
		Width:  btnWidth,
		Height: btnHeight,
	}
	exitRect := rl.Rectangle{
		X:      float32(w)/2 - btnWidth/2,
		Y:      float32(h)/2 + 60 + btnHeight + spacing,
		Width:  btnWidth,
		Height: btnHeight,
	}

	startColor := rl.Maroon
	if blank {
		startColor = rl.DarkGray
	}
	if drawButton(startRect, "Start Game", startColor, rl.Red, rl.White) || tryStart {
		if blank {
			showDiffError = true
		} else {
			showDiffError = false
			showInitMenu = false
			intermission = true
			// rl.PlaySound(sound.ThreeTwoOne)
			initNewGame(
				bulletManager,
				blastManager,
				structureManager,
				player,
				enemyManager,
				stageManager,
			)
			stageManager.CreateNewStage(player.Position)
		}
	}

	if showDiffError {
		errText := "enter a stage number first"
		errSize := int32(18)
		errWidth := rl.MeasureText(errText, errSize)
		rl.DrawText(errText, int32(w)/2-errWidth/2, int32(startRect.Y)-30, errSize, rl.Red)
	}

	exitClicked := drawButton(exitRect, "Exit Game", rl.DarkGray, rl.Gray, rl.White)

	rl.EndDrawing()
	return exitClicked
}

func doGameWon(w, h int) bool {
	if rl.IsCursorHidden() {
		rl.EnableCursor()
	}
	if rl.IsSoundPlaying(sound.Track) {
		rl.StopSound(sound.Track)
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	winText := "THANKS FOR PLAYING! YOU'VE FINISHED THE GAME!"
	fontSizeWin := int32(50)
	textWidthWin := rl.MeasureText(winText, fontSizeWin)
	textX := float32(w)/2 - float32(textWidthWin)/2
	textY := float32(h)/2 - 100
	rl.DrawText(winText, int32(textX), int32(textY), fontSizeWin, rl.Gold)

	btnRect := rl.Rectangle{
		X:      float32(w)/2 - btnWidth/2,
		Y:      float32(h)/2 + 20,
		Width:  btnWidth,
		Height: btnHeight,
	}

	exitClicked := false
	if drawButton(btnRect, "Exit Game", rl.DarkGray, rl.Gray, rl.White) {
		exitClicked = true
	}

	rl.EndDrawing()

	return exitClicked
}

func drawEnemyCount(w int, alive int) {
	text := fmt.Sprintf("ENEMIES: %d", alive)
	fontSize := int32(22)
	textWidth := rl.MeasureText(text, fontSize)
	margin := int32(20)
	rl.DrawText(text, int32(w)/2-textWidth/2, margin, fontSize, rl.NewColor(220, 60, 60, 220))
}

func drawButton(rect rl.Rectangle, text string, baseColor, hoverColor, textColor rl.Color) bool {
	mousePoint := rl.GetMousePosition()
	isHovered := rl.CheckCollisionPointRec(mousePoint, rect)

	currentColor := baseColor

	if isHovered {
		currentColor = hoverColor

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			currentColor = rl.Maroon
		}
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			return true
		}
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
