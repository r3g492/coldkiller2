package main

import (
	"coldkiller2/bullet"
	"coldkiller2/enemy"
	"coldkiller2/input"
	"coldkiller2/killer"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var lastLog = time.Now()

func main() {
	rl.InitWindow(1600, 900, "my new game")
	defer rl.CloseWindow()
	rl.InitAudioDevice()
	rl.SetTargetFPS(144)
	keyMap := input.DefaultWASD()
	bm := bullet.CreateManager()

	em := enemy.CreateManager()
	defer em.Unload()

	p := killer.Init()
	defer p.Unload()

	em.Init()

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyR) {
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
		var ebc = em.Mutate(dt, p)
		em.ProcessAnimation(dt)
		rl.BeginMode3D(p.Camera)
		em.DrawEnemies3D()
		rl.EndMode3D()

		// bullet
		bm.KillerBulletCreate(bc)
		bm.EnemyBulletCreate(ebc)
		bm.Mutate(dt, p, em.Enemies)
		rl.BeginMode3D(p.Camera)
		bm.DrawBullets3D()
		rl.EndMode3D()

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
