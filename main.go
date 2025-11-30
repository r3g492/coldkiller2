package main

import (
	"coldkiller2/bullet"
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
	rl.SetTargetFPS(60)
	p := killer.Init()
	defer p.Unload()
	keyMap := input.DefaultWASD()
	bm := bullet.Manager{}
	for !rl.WindowShouldClose() {
		// seconds
		dt := rl.GetFrameTime()
		mouseLocation := rl.GetMousePosition()
		log(mouseLocation, dt, p)
		bm.Mutate(dt)
		p.Mutate(input.ReadInput(keyMap), dt, &bm)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode3D(p.Camera)
		rl.DrawGrid(1000, 1)
		p.Draw3D()
		bm.DrawBullets()
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
