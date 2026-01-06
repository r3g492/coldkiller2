package main

import (
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
	curr := rl.GetCurrentMonitor()
	w := rl.GetMonitorWidth(curr)
	h := rl.GetMonitorHeight(curr)
	rl.SetWindowSize(w, h)
	rl.ToggleBorderlessWindowed()
	defer rl.CloseWindow()
	rl.DisableCursor()

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

	showMenu := true
	lost := false
	for !rl.WindowShouldClose() {
		if showMenu && lost {
			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)
			rl.DrawText(
				"You Lost! Press R to Restart",
				int32(w/2-200),
				int32(h/2),
				100,
				rl.Red,
			)
			rl.EndDrawing()

			if rl.IsKeyPressed(rl.KeyR) {
				showMenu = false
				lost = true
			}

			continue
		}

		if showMenu {
			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)
			rl.DrawText(
				"Press R to Start",
				int32(w/2-200),
				int32(h/2),
				100,
				rl.Red,
			)
			rl.EndDrawing()

			if rl.IsKeyPressed(rl.KeyR) {
				showMenu = false
			}

			continue
		}

		if gameEnd(p) {
			showMenu = true
			lost = true
			p = resetGame(em, p, bm)
		}

		if rl.IsKeyPressed(rl.KeyEscape) {
			showMenu = true
		}

		if rl.IsKeyPressed(rl.KeyR) {
			p = resetGame(em, p, bm)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// seconds
		dt := rl.GetFrameTime()
		mouseLocation := rl.GetMousePosition()
		log(mouseLocation, dt, p)
		ip := input.ReadInput(keyMap)

		// enemy
		rl.BeginMode3D(p.Camera)
		var ebc = em.Mutate(dt, p)
		em.ProcessAnimation(dt)
		em.DrawEnemies3D(p)
		rl.EndMode3D()

		// player
		bc := p.Mutate(ip, dt, em.GetEnemyBoundingBoxes())
		p.ResolveAnimation()
		p.PlanAnimate(dt)
		p.Animate()
		rl.BeginMode3D(p.Camera)
		p.Draw3D()
		rl.EndMode3D()

		// bullet
		bm.KillerBulletCreate(bc)
		bm.EnemyBulletCreate(ebc)
		bm.Mutate(dt, p, em.Enemies)
		rl.BeginMode3D(p.Camera)
		bm.DrawBullets3D()
		rl.EndMode3D()

		drawCursor(mouseLocation, p)
		rl.DrawText(strconv.Itoa(em.EnemyGenerationLevel), 500, 500, 30, rl.Purple)
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
	rl.DrawRay(rl.NewRay(player.Position, player.TargetDirection), rl.Green)
	rl.EndMode3D()
	rl.DrawCircle(int32(mouseLocation.X), int32(mouseLocation.Y), 5, rl.Green)
}

func gameEnd(player *killer.Killer) bool {
	return player.Health <= 0 && player.ActionTimeLeft <= 0
}
