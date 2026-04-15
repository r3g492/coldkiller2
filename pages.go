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

const splashDuration = 2.0

// doSplash draws the splash screen and returns true when it should end.
func doSplash(timer *float32, tex rl.Texture2D, dt float32, w, h int) bool {
	*timer += dt

	// fade in over first 0.5s, fully visible until 1.5s, fade out over last 0.5s
	var overlay uint8
	t := *timer
	switch {
	case t < 0.5:
		overlay = uint8(255 - t/0.5*255)
	case t < 1.5:
		overlay = 0
	default:
		fade := (t - 1.5) / 0.5
		if fade > 1 {
			fade = 1
		}
		overlay = uint8(fade * 255)
	}

	beginFrame()
	rl.ClearBackground(rl.Black)

	imgW := float32(tex.Width)
	imgH := float32(tex.Height)
	destX := float32(w)/2 - imgW/2
	destY := float32(h)/2 - imgH/2
	rl.DrawTextureV(tex, rl.Vector2{X: destX, Y: destY}, rl.White)

	// gradient: top dark → transparent → bottom dark
	gradH := float32(h)
	for i := 0; i < h; i++ {
		edge := float32(i) / gradH
		var alpha uint8
		if edge < 0.5 {
			alpha = uint8((1 - edge*2) * 180)
		} else {
			alpha = uint8((edge*2 - 1) * 180)
		}
		rl.DrawLine(0, int32(i), int32(w), int32(i), rl.NewColor(0, 0, 0, alpha))
	}

	rl.DrawRectangle(0, 0, int32(w), int32(h), rl.NewColor(0, 0, 0, overlay))

	endFrame()

	return *timer >= splashDuration
}

func doPauseMenu(w, h int, ip input.Input, keyMap input.KeyMap) pauseAction {
	beginFrame()
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

	if rl.IsCursorHidden() {
		rl.EnableCursor()
		scale, offsetX, offsetY := letterbox()
		cx := int32((resumeRect.X+resumeRect.Width/2)*scale + offsetX)
		cy := int32((resumeRect.Y+resumeRect.Height/2)*scale + offsetY)
		rl.SetMousePosition(int(cx), int(cy))
	}

	result := pauseNone
	if drawButton(resumeRect, "Resume", rl.DarkGray, rl.Gray, rl.White) {
		result = pauseResume
	}
	if drawButton(quitRect, "Quit to Menu", rl.DarkGray, rl.Gray, rl.White) {
		result = pauseQuitToMenu
	}

	drawInputOverlay(w, h, ip, keyMap, true)

	endFrame()
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
	keyMap input.KeyMap,
) {
	const totalLoadSteps = 9

	beginFrame()
	rl.ClearBackground(rl.DarkGray)

	diffInfoText := fmt.Sprintf("%d / %d", stageManager.Difficulty, len(stage.Stages))
	diffInfoSize := int32(40)
	diffInfoWidth := rl.MeasureText(diffInfoText, diffInfoSize)
	rl.DrawText(diffInfoText, int32(w)/2-diffInfoWidth/2, int32(h)/2-40, diffInfoSize, rl.RayWhite)

	intermissionTimer += dt

	ready := intermissionLoadStep >= totalLoadSteps

	if !ready {
		// loading in progress: show bar and execute one init step per frame

		if intermissionLoadStep < totalLoadSteps {
			switch intermissionLoadStep {
			case 0:
				stageManager.GenerateNewStage()
			case 1:
				bulletManager.Init()
			case 2:
				blastManager.Init()
			case 3:
				structureManager.Init()
			case 4:
				player.Init()
			case 5:
				enemyManager.Init(player)
			case 6:
				stageManager.Init(structureManager, enemyManager, player)
			case 7:
				stage.InitStages()
			case 8:
				stageManager.CreateNewStage(player.Position)
			}
			intermissionLoadStep++
		}
	} else {
		// loading done: show start button
		startRect := rl.Rectangle{
			X:      float32(w)/2 - btnWidth/2,
			Y:      float32(h)/2 + 35,
			Width:  btnWidth,
			Height: btnHeight,
		}
		if rl.IsCursorHidden() {
			rl.EnableCursor()
			scale, offsetX, offsetY := letterbox()
			cx := int32((startRect.X+startRect.Width/2)*scale + offsetX)
			cy := int32((startRect.Y+startRect.Height/2)*scale + offsetY)
			rl.SetMousePosition(int(cx), int(cy))
		}
		if drawButton(startRect, "Go", rl.Maroon, rl.Red, rl.White) || rl.IsKeyPressed(rl.KeyEnter) {
			intermission = false
			intermissionTimer = 0
			intermissionLoadStep = 0
		}
	}

	endFrame()
}

func doInitMenu(
	stageManager *stage.Manager,
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

	beginFrame()
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
		}
	}

	if showDiffError {
		errText := "enter a stage number first"
		errSize := int32(18)
		errWidth := rl.MeasureText(errText, errSize)
		rl.DrawText(errText, int32(w)/2-errWidth/2, int32(startRect.Y)-30, errSize, rl.Red)
	}

	exitClicked := drawButton(exitRect, "Exit Game", rl.DarkGray, rl.Gray, rl.White)

	endFrame()
	return exitClicked
}

func doGameWon(w, h int) bool {
	if rl.IsCursorHidden() {
		rl.EnableCursor()
	}
	if rl.IsSoundPlaying(sound.Track) {
		rl.StopSound(sound.Track)
	}

	beginFrame()
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

	endFrame()

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
	mousePoint := virtualMousePosition()
	isHovered := rl.CheckCollisionPointRec(mousePoint, rect)

	currentColor := baseColor

	if isHovered {
		currentColor = hoverColor

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			currentColor = rl.Maroon
		}
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
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

func drawInputOverlay(w, h int, ip input.Input, keyMap input.KeyMap, midTop bool) {
	const (
		keySize      = 45
		spacing      = 6
		fontSize     = 18
		descSize     = 9
		rightMargin  = 40
		bottomMargin = 90
		topMargin    = 240
	)

	labelUp := input.GetKeyName(keyMap.Up)
	labelLeft := input.GetKeyName(keyMap.Left)
	labelDown := input.GetKeyName(keyMap.Down)
	labelRight := input.GetKeyName(keyMap.Right)
	labelReload := input.GetKeyName(keyMap.Reload)
	labelEnd := input.GetKeyName(keyMap.EndGame)

	totalWidth := (keySize * 8) + (spacing * 7) + 50
	var baseX, baseY float32
	if midTop {
		baseX = float32(w)/2 - float32(totalWidth)/2
		baseY = float32(topMargin)
	} else {
		baseX = float32(w) - float32(totalWidth) - rightMargin
		baseY = float32(h) - (keySize * 2) - spacing - bottomMargin
	}

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

	drawKey(currX, baseY+keySize+spacing, keySize+30, "LMB", "SHOOT", rl.IsMouseButtonDown(keyMap.Fire))
	currX += keySize + 30 + spacing

	drawKey(currX, baseY+keySize+spacing, keySize+30, "SPC", "DASH", ip.DashPressed)
}

func drawCursor(mouseLocation rl.Vector2, player *killer.Killer) {
	mouseRay := rl.GetScreenToWorldRayEx(mouseLocation, player.Camera, VirtualWidth, VirtualHeight)

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
