package main

import (
	"coldkiller2/background"
	"coldkiller2/bullet"
	"coldkiller2/enemy"
	"coldkiller2/input"
	"coldkiller2/killer"
	"coldkiller2/sound"
	"fmt"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var lastLog = time.Now()

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
		// seconds
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
			rl.DrawText("YOU LOSE", int32(w/2-150), int32(h/2-150), 60, rl.Red)

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

		drawCursor(mouseLocation, p)
		rl.DrawText(strconv.Itoa(em.EnemyGenerationLevel), 500, 500, 30, rl.Purple)
		rl.DrawText(strconv.Itoa(bm.PlayerXp), 700, 500, 30, rl.Red)
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
