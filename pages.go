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
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type initMenuAction int

const (
	initMenuNone initMenuAction = iota
	initMenuContinue
	initMenuStart
	initMenuOptions
	initMenuExit
	initMenuShowCleared
	initMenuTestStart
	initMenuEditor
)

type upgradeInfo struct {
	name    string
	current func(*killer.Killer) string
	next    func(*killer.Killer) string
	capped  func(*killer.Killer) bool
	apply   func(*killer.Killer)
}

const (
	slowDurationStep = 0.5
	slowDurationMax  = 15

	ammoCapacityStep = 3
	ammoCapacityMax  = 1000

	moveSpeedStep = 0.25
	moveSpeedMax  = 12
)

var allUpgrades = []upgradeInfo{
	{
		name:    "Slow Duration",
		current: func(p *killer.Killer) string { return fmt.Sprintf("%.1fs", p.SlowTimeDuration) },
		next:    func(p *killer.Killer) string { return fmt.Sprintf("%.1fs", p.SlowTimeDuration+slowDurationStep) },
		capped:  func(p *killer.Killer) bool { return p.SlowTimeDuration >= slowDurationMax },
		apply:   func(p *killer.Killer) { p.SlowTimeDuration += slowDurationStep },
	},
	{
		name:    "Ammo Capacity",
		current: func(p *killer.Killer) string { return fmt.Sprintf("%d", p.AmmoCapacity) },
		next:    func(p *killer.Killer) string { return fmt.Sprintf("%d", p.AmmoCapacity+ammoCapacityStep) },
		capped:  func(p *killer.Killer) bool { return p.AmmoCapacity >= ammoCapacityMax },
		apply:   func(p *killer.Killer) { p.AmmoCapacity += ammoCapacityStep; p.Ammo = p.AmmoCapacity },
	},
	{
		name:    "Move Speed",
		current: func(p *killer.Killer) string { return fmt.Sprintf("%.2f", p.MoveSpeed) },
		next:    func(p *killer.Killer) string { return fmt.Sprintf("%.2f", p.MoveSpeed+moveSpeedStep) },
		capped:  func(p *killer.Killer) bool { return p.MoveSpeed >= moveSpeedMax },
		apply:   func(p *killer.Killer) { p.MoveSpeed += moveSpeedStep },
	},
}

func drawFullyUpgradedButton(rect rl.Rectangle, name string) {
	dim := rl.NewColor(60, 60, 60, 255)
	rl.DrawRectangleRec(rect, rl.Fade(dim, 0.3))
	rl.DrawRectangleLinesEx(rect, 3, dim)

	const nameSize = int32(13)
	const labelSize = int32(12)

	nw := rl.MeasureText(name, nameSize)
	rl.DrawText(name, int32(rect.X+rect.Width/2)-nw/2, int32(rect.Y)+10, nameSize, rl.Gray)

	label := "Fully Upgraded"
	lw := rl.MeasureText(label, labelSize)
	rl.DrawText(label, int32(rect.X+rect.Width/2)-lw/2, int32(rect.Y)+int32(rect.Height)-labelSize-10, labelSize, rl.DarkGray)
}

func drawUpgradeButton(rect rl.Rectangle, p *killer.Killer, u *upgradeInfo, baseColor, hoverColor rl.Color) bool {
	mousePoint := virtualMousePosition()
	isHovered := rl.CheckCollisionPointRec(mousePoint, rect)
	color := baseColor
	if isHovered {
		color = hoverColor
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			color = rl.Maroon
		}
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			return true
		}
	}
	rl.DrawRectangleRec(rect, rl.Fade(color, 0.3))
	rl.DrawRectangleLinesEx(rect, 3, color)

	const nameSize = int32(15)
	const valSize = int32(14)

	nw := rl.MeasureText(u.name, nameSize)
	rl.DrawText(u.name, int32(rect.X+rect.Width/2)-nw/2, int32(rect.Y)+10, nameSize, rl.White)

	valText := u.current(p) + " -> " + u.next(p)
	vw := rl.MeasureText(valText, valSize)
	rl.DrawText(valText, int32(rect.X+rect.Width/2)-vw/2, int32(rect.Y)+int32(rect.Height)-valSize-10, valSize, rl.LightGray)

	return false
}

var (
	btnWidth  = float32(400)
	btnHeight = float32(80)
	spacing   = float32(20)
)

type pauseAction int

const (
	pauseNone pauseAction = iota
	pauseResume
	pauseQuitToMenu
)

const splashDuration = 2.0

func doSplash(timer *float32, tex rl.Texture2D, dt float32, w, h int) bool {
	*timer += dt

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

func doPauseMenu(w, h int, ip input.Input, keyMap input.KeyMap, slowAvailable bool) pauseAction {
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
	drawInputOverlay(w, h, ip, keyMap, slowAvailable)
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
) bool {
	const totalLoadSteps = 8

	gaveUp := false

	beginFrame()
	rl.ClearBackground(rl.DarkGray)

	diffInfoText := fmt.Sprintf("%d / %d", stageManager.Difficulty, len(stage.Stages))
	diffInfoSize := int32(40)
	diffInfoWidth := rl.MeasureText(diffInfoText, diffInfoSize)
	rl.DrawText(diffInfoText, int32(w)/2-diffInfoWidth/2, int32(h)/2-40, diffInfoSize, rl.RayWhite)

	drawBestKPS(w, 60)
	drawLatestKPS(w, 90)

	intermissionTimer += dt

	ready := intermissionLoadStep >= totalLoadSteps

	if !ready {
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
			stageManager.CreateNewStage(player.Position)
		}
		intermissionLoadStep++
	} else if stageManager.Difficulty == 1 {
		player.RefillSlowTime()
		intermission = false
		intermissionTimer = 0
		intermissionLoadStep = 0
	} else {
		halfWidth := btnWidth/2 - spacing/2
		startX := float32(w)/2 - btnWidth/2
		startY := float32(h)/2 + 35

		if rl.IsCursorHidden() {
			rl.EnableCursor()
			scale, offsetX, offsetY := letterbox()
			cx := int32((startX+halfWidth/2)*scale + offsetX)
			cy := int32((startY+btnHeight/2)*scale + offsetY)
			rl.SetMousePosition(int(cx), int(cy))
		}

		var chosen *upgradeInfo
		for i := range allUpgrades {
			u := &allUpgrades[i]
			col := i % 2
			row := i / 2
			rect := rl.Rectangle{
				X:      startX + float32(col)*(halfWidth+spacing),
				Y:      startY + float32(row)*(btnHeight+spacing),
				Width:  halfWidth,
				Height: btnHeight,
			}
			if u.capped(player) {
				drawFullyUpgradedButton(rect, u.name)
				continue
			}
			if drawUpgradeButton(rect, player, u, rl.Maroon, rl.Red) {
				chosen = u
			}
		}

		if chosen != nil {
			chosen.apply(player)
			player.RefillSlowTime()
			intermission = false
			intermissionTimer = 0
			intermissionLoadStep = 0
		}
	}
	endFrame()
	return gaveUp
}

func doInitMenu(w, h int) initMenuAction {
	rl.StopSound(sound.Track)

	if rl.IsCursorHidden() {
		rl.EnableCursor()
	}

	beginFrame()
	rl.ClearBackground(rl.Black)

	title := "Kill Per Second"
	titleSize := int32(90)
	titleW := rl.MeasureText(title, titleSize)
	t := float32(rl.GetTime())
	hue := t * 60
	hue -= float32(int(hue/360)) * 360
	titleColor := rl.ColorFromHSV(hue, 0.85, 1.0)
	titleX := int32(w)/2 - titleW/2
	titleY := int32(h)/2 - 280
	rl.DrawText(title, titleX, titleY, titleSize, titleColor)

	const menuSpacing = float32(100)
	menuStartY := float32(h)/2 - 1.5*menuSpacing

	makeMenuRect := func(i int) rl.Rectangle {
		return rl.Rectangle{
			X:      float32(w)/2 - btnWidth/2,
			Y:      menuStartY + float32(i)*menuSpacing,
			Width:  btnWidth,
			Height: btnHeight,
		}
	}

	resumable := currentConfig.StageLevel >= 1 && currentConfig.StageLevel < len(stage.Stages)

	result := initMenuNone
	switch {
	case resumable:
		label := fmt.Sprintf("Retry (Stage %d)", currentConfig.StageLevel+1)
		if drawButton(makeMenuRect(0), label, rl.NewColor(0, 90, 40, 255), rl.NewColor(0, 150, 70, 255), rl.White) {
			result = initMenuContinue
		}
	case currentConfig.StageLevel >= len(stage.Stages):
		if drawClearedButton(makeMenuRect(0), "GAME CLEARED") {
			result = initMenuShowCleared
		}
	default:
		drawDisabledButton(makeMenuRect(0), "Retry")
	}
	if drawButton(makeMenuRect(1), "New Game", rl.Maroon, rl.Red, rl.White) || rl.IsKeyPressed(rl.KeyEnter) {
		result = initMenuStart
	}
	if drawButton(makeMenuRect(2), "Options", rl.NewColor(50, 50, 80, 255), rl.NewColor(80, 80, 130, 255), rl.White) {
		result = initMenuOptions
	}
	if drawButton(makeMenuRect(3), "Exit", rl.DarkGray, rl.Gray, rl.White) {
		result = initMenuExit
	}

	if testMode {
		if a := drawTestPanel(w, h); a != initMenuNone {
			result = a
		}
	}

	endFrame()
	return result
}

// drawTestPanel renders the test-config editor in the empty space to the right
// of the centered main menu. It mutates testConfig in place and returns the
// action triggered by its buttons (or initMenuNone). Only shown while testMode
// is on.
func drawTestPanel(w, h int) initMenuAction {
	panelX := float32(w)/2 + btnWidth/2 + 70
	startY := float32(h)/2 - 150
	const rowH = float32(64)

	rl.DrawText("TEST CONFIG", int32(panelX), int32(startY)-50, 30, rl.Gold)

	if dec, inc := drawTestStepper(panelX, startY+0*rowH, "Slow Dur", fmt.Sprintf("%.2f", testConfig.SlowDuration)); dec || inc {
		testConfig.SlowDuration += stepDelta(dec, inc, 0.5)
		if testConfig.SlowDuration < 0 {
			testConfig.SlowDuration = 0
		}
	}
	if dec, inc := drawTestStepper(panelX, startY+1*rowH, "Ammo Cap", fmt.Sprintf("%d", testConfig.AmmoCapacity)); dec || inc {
		testConfig.AmmoCapacity += int32(stepDelta(dec, inc, 5))
		if testConfig.AmmoCapacity < 1 {
			testConfig.AmmoCapacity = 1
		}
	}
	if dec, inc := drawTestStepper(panelX, startY+2*rowH, "Move Spd", fmt.Sprintf("%.1f", testConfig.MoveSpeed)); dec || inc {
		testConfig.MoveSpeed += stepDelta(dec, inc, 0.5)
		if testConfig.MoveSpeed < 0.5 {
			testConfig.MoveSpeed = 0.5
		}
	}
	if dec, inc := drawTestStepper(panelX, startY+3*rowH, "Stage", fmt.Sprintf("%d / %d", testConfig.Stage, len(stage.Stages))); dec || inc {
		testConfig.Stage += int(stepDelta(dec, inc, 1))
		if testConfig.Stage < 1 {
			testConfig.Stage = 1
		}
		if testConfig.Stage > len(stage.Stages) {
			testConfig.Stage = len(stage.Stages)
		}
	}

	startRect := rl.Rectangle{X: panelX, Y: startY + 4*rowH + 10, Width: btnWidth, Height: btnHeight}
	if drawButton(startRect, "Test Start", rl.NewColor(120, 80, 0, 255), rl.NewColor(180, 120, 0, 255), rl.White) {
		return initMenuTestStart
	}

	editorRect := rl.Rectangle{X: panelX, Y: startY + 4*rowH + 10 + btnHeight + spacing, Width: btnWidth, Height: btnHeight}
	if drawButton(editorRect, "Stage Editor", rl.NewColor(40, 70, 110, 255), rl.NewColor(70, 110, 170, 255), rl.White) {
		return initMenuEditor
	}

	return initMenuNone
}

// stepDelta returns +mag for an increment click, -mag for a decrement click.
func stepDelta(dec, inc bool, mag float32) float32 {
	if inc {
		return mag
	}
	if dec {
		return -mag
	}
	return 0
}

// drawTestStepper draws a "label [-] value [+]" row and reports which arrow was
// clicked this frame.
func drawTestStepper(x, y float32, label, valueText string) (dec, inc bool) {
	rl.DrawText(label, int32(x), int32(y)+10, 22, rl.RayWhite)

	const btnSize = float32(40)
	minusRect := rl.Rectangle{X: x + 230, Y: y, Width: btnSize, Height: btnSize}
	plusRect := rl.Rectangle{X: x + 380, Y: y, Width: btnSize, Height: btnSize}
	dec = drawButtonSized(minusRect, "-", 28, rl.DarkGray, rl.Gray, rl.White)
	inc = drawButtonSized(plusRect, "+", 28, rl.DarkGray, rl.Gray, rl.White)

	const valSize = int32(22)
	vw := rl.MeasureText(valueText, valSize)
	rl.DrawText(valueText, int32(x+325)-vw/2, int32(y)+10, valSize, rl.SkyBlue)
	return
}

func doGameWon(w, h int, buttonLabel string, winTex rl.Texture2D) bool {
	if rl.IsCursorHidden() {
		rl.EnableCursor()
	}
	if rl.IsSoundPlaying(sound.Track) {
		rl.StopSound(sound.Track)
	}

	beginFrame()
	rl.ClearBackground(rl.Black)

	const (
		winNativeW = float32(1000)
		winNativeH = float32(1200)
		imgH       = float32(540)
	)
	imgW := imgH * winNativeW / winNativeH
	src := rl.Rectangle{X: 0, Y: 0, Width: float32(winTex.Width), Height: float32(winTex.Height)}
	dst := rl.Rectangle{X: float32(w)/2 - imgW/2, Y: 30, Width: imgW, Height: imgH}
	rl.DrawTexturePro(winTex, src, dst, rl.Vector2{}, 0, rl.White)

	winText := "THANKS FOR PLAYING! YOU'VE FINISHED THE GAME!"
	fontSizeWin := int32(40)
	textWidthWin := rl.MeasureText(winText, fontSizeWin)
	textX := float32(w)/2 - float32(textWidthWin)/2
	textY := dst.Y + dst.Height + 20
	rl.DrawText(winText, int32(textX), int32(textY), fontSizeWin, rl.Gold)

	drawBestKPS(w, int32(textY)+fontSizeWin+15)

	btnRect := rl.Rectangle{
		X:      float32(w)/2 - btnWidth/2,
		Y:      textY + float32(fontSizeWin) + 55,
		Width:  btnWidth,
		Height: btnHeight,
	}

	clicked := drawButton(btnRect, buttonLabel, rl.DarkGray, rl.Gray, rl.White)

	creditLines := []string{
		"Thanks to:",
		"CaveTroll for the artwork",
		"Sounds: 1911-reload-6248.mp3, data_pion-st1-footstep-sfx-323053.mp3, shotgun-03-38220.mp3",
	}
	creditFontSize := int32(20)
	creditY := int32(btnRect.Y+btnHeight) + 40
	for _, line := range creditLines {
		lineWidth := rl.MeasureText(line, creditFontSize)
		lineX := int32(float32(w)/2 - float32(lineWidth)/2)
		rl.DrawText(line, lineX, creditY, creditFontSize, rl.LightGray)
		creditY += creditFontSize + 8
	}

	endFrame()
	return clicked
}

func drawKPS(w int, y int32) {
	kps := float32(0)
	if kpsTimeSec > 0 {
		kps = float32(kpsKills) / kpsTimeSec
	}
	drawKpsText(w, y, kps, "Kills Per Second")
}

func drawBestKPS(w int, y int32) {
	drawKpsText(w, y, kpsBest, "Best Kills Per Second")
}

func drawLatestKPS(w int, y int32) {
	drawKpsText(w, y, kpsLatest, "Latest Kills Per Second")
}

func drawKpsText(w int, y int32, kps float32, label string) {
	text := fmt.Sprintf("%s: %.2f", label, kps)
	fontSize := int32(22)
	textWidth := rl.MeasureText(text, fontSize)
	baseX := int32(w)/2 - textWidth/2

	if kps <= 2 {
		rl.DrawText(text, baseX, y, fontSize, rl.NewColor(220, 60, 60, 220))
		return
	}

	t := float32(rl.GetTime())
	hue := t * 180
	hue -= float32(int(hue/360)) * 360
	color := rl.ColorFromHSV(hue, 1.0, 1.0)

	shake := float32(0)
	glow := 0
	if kps > 3 {
		shake = 1.5
		glow = 3
	}
	if kps > 4 {
		shake = 3.0
		glow = 6
	}

	var dx, dy int32
	if shake > 0 {
		dx = int32((rand.Float32()*2 - 1) * shake)
		dy = int32((rand.Float32()*2 - 1) * shake)
	}

	for i := 1; i <= glow; i++ {
		alpha := uint8(120 - i*15)
		if alpha < 30 {
			alpha = 30
		}
		glowColor := color
		glowColor.A = alpha
		offset := int32(i)
		rl.DrawText(text, baseX+dx-offset, y+dy, fontSize, glowColor)
		rl.DrawText(text, baseX+dx+offset, y+dy, fontSize, glowColor)
		rl.DrawText(text, baseX+dx, y+dy-offset, fontSize, glowColor)
		rl.DrawText(text, baseX+dx, y+dy+offset, fontSize, glowColor)
	}

	rl.DrawText(text, baseX+dx, y+dy, fontSize, color)
}

func drawButton(rect rl.Rectangle, text string, baseColor, hoverColor, textColor rl.Color) bool {
	return drawButtonSized(rect, text, 30, baseColor, hoverColor, textColor)
}

func drawDisabledButton(rect rl.Rectangle, text string) {
	dim := rl.NewColor(60, 60, 60, 255)
	rl.DrawRectangleRec(rect, rl.Fade(dim, 0.2))
	rl.DrawRectangleLinesEx(rect, 3, dim)

	const fontSize = int32(30)
	textWidth := rl.MeasureText(text, fontSize)
	textX := int32(rect.X + (rect.Width / 2) - float32(textWidth/2))
	textY := int32(rect.Y + (rect.Height / 2) - float32(fontSize/2))
	rl.DrawText(text, textX, textY, fontSize, rl.Gray)
}

func drawClearedButton(rect rl.Rectangle, text string) bool {
	pulse := 0.7 + 0.3*float32(math.Abs(math.Sin(rl.GetTime()*2)))
	gold := rl.Gold

	hovered := rl.CheckCollisionPointRec(virtualMousePosition(), rect)
	fillAlpha := float32(0.15)
	if hovered {
		fillAlpha = 0.3
	}

	rl.DrawRectangleRec(rect, rl.Fade(gold, fillAlpha*pulse))
	rl.DrawRectangleLinesEx(rect, 3, rl.Fade(gold, pulse))

	const fontSize = int32(30)
	textWidth := rl.MeasureText(text, fontSize)
	textX := int32(rect.X + (rect.Width / 2) - float32(textWidth/2))
	textY := int32(rect.Y + (rect.Height / 2) - float32(fontSize/2))
	rl.DrawText(text, textX, textY, fontSize, gold)

	return hovered && rl.IsMouseButtonReleased(rl.MouseLeftButton)
}

func drawButtonSized(rect rl.Rectangle, text string, fontSize int32, baseColor, hoverColor, textColor rl.Color) bool {
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

	textWidth := rl.MeasureText(text, fontSize)
	textX := int32(rect.X + (rect.Width / 2) - float32(textWidth/2))
	textY := int32(rect.Y + (rect.Height / 2) - float32(fontSize/2))

	rl.DrawText(text, textX, textY, fontSize, textColor)

	return false
}

func drawKey(x, y float32, width float32, label string, desc string, active, disabled bool) {
	const (
		keySize    = 45
		fontSize   = 18
		descSize   = 10
		borderSize = 2
	)

	rect := rl.Rectangle{X: x, Y: y, Width: width, Height: keySize}

	bgColor := rl.NewColor(50, 50, 60, 210)
	borderCol := rl.Fade(rl.LightGray, 0.85)
	textCol := rl.RayWhite
	descCol := rl.Fade(rl.RayWhite, 0.85)

	if active {
		bgColor = rl.NewColor(230, 41, 55, 220)
		borderCol = rl.Red
		textCol = rl.White
		descCol = rl.Fade(rl.White, 0.9)
	} else if disabled {
		bgColor = rl.NewColor(35, 35, 40, 150)
		borderCol = rl.Fade(rl.Gray, 0.5)
		textCol = rl.Fade(rl.Gray, 0.7)
		descCol = rl.Fade(rl.Gray, 0.55)
	}

	rl.DrawRectangleRec(rect, bgColor)
	rl.DrawRectangleLinesEx(rect, borderSize, borderCol)

	tw := rl.MeasureText(label, fontSize)
	rl.DrawText(label, int32(x+width/2)-tw/2, int32(y+keySize/2)-fontSize/2-4, fontSize, textCol)

	dtw := rl.MeasureText(desc, descSize)
	rl.DrawText(desc, int32(x+width/2)-dtw/2, int32(y+keySize)-descSize-5, descSize, descCol)
}

func drawInputOverlay(w, h int, ip input.Input, keyMap input.KeyMap, slowAvailable bool) {
	const (
		keySize = 45
		spacing = 6
		margin  = 30
		wideKey = keySize + 30
	)

	labelUp := input.GetKeyName(keyMap.Up)
	labelLeft := input.GetKeyName(keyMap.Left)
	labelDown := input.GetKeyName(keyMap.Down)
	labelRight := input.GetKeyName(keyMap.Right)
	labelReload := input.GetKeyName(keyMap.Reload)
	labelEnd := input.GetKeyName(keyMap.EndGame)

	drawKey(margin, margin, keySize, labelEnd, "END", rl.IsKeyDown(keyMap.EndGame), false)

	clusterCols := 4
	clusterWidth := float32(clusterCols*keySize + (clusterCols-1)*spacing)
	clusterX := float32(w)/2 - clusterWidth/2

	spaceRowY := float32(h) - margin - keySize
	midRowY := spaceRowY - keySize - spacing
	topRowY := midRowY - keySize - spacing

	col := func(i int) float32 { return clusterX + float32(i*(keySize+spacing)) }

	drawKey(col(1), topRowY, keySize, labelUp, "UP", ip.MoveUp, false)
	drawKey(col(3), topRowY, keySize, labelReload, "RELOAD", ip.ReloadPressed, false)

	drawKey(col(0), midRowY, keySize, labelLeft, "LEFT", ip.MoveLeft, false)
	drawKey(col(1), midRowY, keySize, labelDown, "DOWN", ip.MoveDown, false)
	drawKey(col(2), midRowY, keySize, labelRight, "RIGHT", ip.MoveRight, false)

	drawKey(clusterX, spaceRowY, clusterWidth, input.GetKeyName(keyMap.SlowTime), "SLOW", ip.SlowTimeDown && slowAvailable, !slowAvailable)

	lmbX := clusterX + clusterWidth + keySize
	drawKey(lmbX, midRowY, wideKey, "LMB", "SHOOT", rl.IsMouseButtonDown(keyMap.Fire), false)
}

func drawEscOverlay() {
	const (
		keySize = 45
		margin  = 30
	)
	drawKey(margin, margin, keySize, "ESC", "pause or keymap", rl.IsKeyDown(rl.KeyEscape), false)
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

	lineEnd := target3D
	aimVec := rl.Vector3Subtract(target3D, player.Position)
	if rl.Vector3Length(aimVec) > player.Range {
		aimDir := rl.Vector3Normalize(aimVec)
		lineEnd = rl.Vector3Add(player.Position, rl.Vector3Scale(aimDir, player.Range))
	}

	rl.BeginMode3D(player.Camera)

	rl.DrawLine3D(player.Position, lineEnd, rl.Green)

	rl.DrawSphere(target3D, 0.1, rl.Green)

	rl.EndMode3D()
}

func drawHitMarker(mouseLocation rl.Vector2) {
	if hitMarkerTimer <= 0 {
		return
	}

	const (
		gap       = 6.0
		length    = 11.0
		thickness = 3.0
	)
	t := hitMarkerTimer / hitMarkerDuration
	alpha := uint8(t * t * 255)
	color := rl.NewColor(255, 60, 60, alpha)
	cx, cy := mouseLocation.X, mouseLocation.Y

	corners := [4][2]float32{{-1, -1}, {1, -1}, {-1, 1}, {1, 1}}
	for _, c := range corners {
		inner := rl.Vector2{X: cx + c[0]*gap, Y: cy + c[1]*gap}
		outer := rl.Vector2{X: cx + c[0]*(gap+length), Y: cy + c[1]*(gap+length)}
		rl.DrawLineEx(inner, outer, thickness, color)
	}
}
