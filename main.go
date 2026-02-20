package main

import (
	"coldkiller2/background"
	"coldkiller2/bullet"
	"coldkiller2/enemy"
	"coldkiller2/input"
	"coldkiller2/killer"
	"coldkiller2/sound"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var lastLog = time.Now()
var lastScore int

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowUndecorated)
	rl.InitWindow(0, 0, "coldkiller2")

	targetMonitor := 1
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
	rl.SetTargetFPS(144)
	keyMap := input.DefaultWASD()
	bm := bullet.CreateManager()
	defer bm.Unload()

	em := enemy.CreateManager()
	defer em.Unload()

	p := killer.Init()
	defer p.Unload()

	em.Init(p)
	sound.Init()

	rl.DisableCursor()

	showMenu := true
	lost := false

	btnWidth, btnHeight := float32(400), float32(80)
	btnRect := rl.Rectangle{
		X:      float32(w)/2 - btnWidth/2,
		Y:      float32(h)/2 - btnHeight/2,
		Width:  btnWidth,
		Height: btnHeight,
	}

	background.InitEnvironment()

	for !rl.WindowShouldClose() {
		if rl.IsCursorHidden() {
			rl.EnableCursor()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		buttonText := "Start Game"
		if drawButton(btnRect, buttonText) {
			showMenu = false
			break
		}

		rl.EndDrawing()
		continue
	}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		mouseLocation := rl.GetMousePosition()
		log(mouseLocation, dt, p)
		ip := input.ReadInput(keyMap)

		if showMenu {
			if rl.IsCursorHidden() {
				rl.EnableCursor()
			}

			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)

			buttonText := "Restart Game"
			xpText := fmt.Sprintf("%d KILL", lastScore)

			fontSize := int32(60)
			textWidth := rl.MeasureText(xpText, fontSize)
			rl.DrawText(xpText, int32(w)/2-textWidth/2, int32(h/2-300), fontSize, rl.Red)

			if drawButton(btnRect, buttonText) || ip.ResetGamePressed {
				showMenu = false
				if lost {
					p = resetGame(em, p, bm)
					lost = false
				}
			}

			rl.EndDrawing()
			continue
		}

		if gameEnd(p) {
			rl.StopSound(sound.Track)
			rl.PlaySound(sound.YouLose)
			showMenu = true
			lost = true
			lastScore = bm.PlayerXp
			p = resetGame(em, p, bm)
		}

		if !rl.IsSoundPlaying(sound.Track) && !showMenu {
			rl.PlaySound(sound.Track)
		}

		if !rl.IsCursorHidden() {
			rl.DisableCursor()
		}

		if ip.EndGamePressed {
			showMenu = true
		}

		if ip.ResetGamePressed {
			p = resetGame(em, p, bm)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(10, 10, 15, 255))
		background.DrawCleanEnvironment(p)
		// enemy
		rl.BeginMode3D(p.Camera)
		var ebc = em.Mutate(dt, p)
		em.ProcessAnimation(dt, p)
		em.DrawEnemies3D(p)
		rl.EndMode3D()
		em.DrawEnemiesUi(p)

		// player
		bc := p.Mutate(ip, dt, em.GetEnemyBoundingBoxes())
		p.ResolveAnimation()
		p.PlanAnimate(dt)
		p.Animate()
		rl.BeginMode3D(p.Camera)
		p.Draw3D()
		rl.EndMode3D()
		p.DrawUI()

		// bullet
		bm.KillerBulletCreate(bc)
		bm.EnemyBulletCreate(ebc)
		bm.Mutate(dt, p, em.Enemies)
		rl.BeginMode3D(p.Camera)
		bm.DrawBullets3D()
		rl.EndMode3D()
		drawInputOverlay(w, h, ip, keyMap)
		drawCursor(mouseLocation, p)
		// rl.DrawText(strconv.Itoa(em.EnemyGenerationLevel), 500, 500, 30, rl.Purple)
		// rl.DrawText(strconv.Itoa(bm.PlayerXp), 700, 500, 30, rl.Red)
		rl.EndDrawing()
	}
}

func resetGame(em *enemy.Manager, p *killer.Killer, bm *bullet.Manager) *killer.Killer {
	em.Unload()
	p.Unload()
	bm.Unload()

	p = killer.Init()
	em.Init(p)
	return p
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
	rl.BeginMode3D(player.Camera)
	// rl.DrawRay(rl.NewRay(player.Position, player.TargetDirection), rl.Green)
	rl.EndMode3D()
	rl.DrawCircle(int32(mouseLocation.X), int32(mouseLocation.Y), 2.5, rl.Green)
}

func gameEnd(player *killer.Killer) bool {
	return player.Health <= 0 && player.ActionTimeLeft <= 0
}

func drawButton(rect rl.Rectangle, text string) bool {
	mousePoint := rl.GetMousePosition()
	isPressed := rl.CheckCollisionPointRec(mousePoint, rect) && (rl.IsMouseButtonDown(rl.MouseLeftButton) || rl.IsMouseButtonReleased(rl.MouseLeftButton))

	color := rl.Red
	if isPressed {
		color = rl.Maroon
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			return true
		}
	}

	rl.DrawRectangleRec(rect, rl.Fade(color, 0.3))
	rl.DrawRectangleLinesEx(rect, 3, color)

	fontSize := int32(30)
	textWidth := rl.MeasureText(text, fontSize)
	textX := int32(rect.X + (rect.Width / 2) - float32(textWidth/2))
	textY := int32(rect.Y + (rect.Height / 2) - float32(fontSize/2))

	rl.DrawText(text, textX, textY, fontSize, color)

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
