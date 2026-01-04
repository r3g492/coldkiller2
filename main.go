package main

import (
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

	em := enemy.CreateManager()
	defer em.Unload()

	p := killer.Init()
	defer p.Unload()

	em.Init()
	sound.Init()

	showMenu := true

	for !rl.WindowShouldClose() {
		if showMenu {
			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)
			rl.DrawText(
				"HELLO",
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

		if rl.IsKeyPressed(rl.KeyEscape) {
			showMenu = true
		}

		if rl.IsKeyPressed(rl.KeyR) {
			em.Unload()
			p.Unload()

			p = killer.Init()
			em.Init()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// seconds
		dt := rl.GetFrameTime()
		mouseLocation := rl.GetMousePosition()
		log(mouseLocation, dt, p)
		ip := input.ReadInput(keyMap)

		// player
		bc := p.Mutate(ip, dt, em.GetEnemyBoundingBoxes())
		p.ResolveAnimation()
		p.PlanAnimate(dt)
		p.Animate()
		rl.BeginMode3D(p.Camera)
		p.Draw3D()
		rl.EndMode3D()

		// enemy
		rl.BeginMode3D(p.Camera)
		var ebc = em.Mutate(dt, p)
		em.ProcessAnimation(dt)
		em.DrawEnemies3D(p)
		rl.EndMode3D()

		// bullet
		bm.KillerBulletCreate(bc)
		bm.EnemyBulletCreate(ebc)
		bm.Mutate(dt, p, em.Enemies)
		rl.BeginMode3D(p.Camera)
		bm.DrawBullets3D()
		rl.EndMode3D()

		drawCursor(mouseLocation, p)

		rl.EndDrawing()
	}
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
