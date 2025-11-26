package main

import (
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
	rl.SetTargetFPS(60)
	player := killer.Init()
	defer player.Unload()
	keyMap := input.DefaultWASD()
	camera3d := player.GetCamera()
	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		mouseLocation := rl.GetMousePosition()
		log(mouseLocation, dt)
		player.Control(input.ReadInput(keyMap), dt)
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode3D(camera3d)
		rl.DrawGrid(1000, 1)
		player.Draw3D()
		rl.EndMode3D()
		rl.EndDrawing()
	}
}

func log(mouseLocation rl.Vector2, dt float32) {
	if time.Since(lastLog) >= 1000*time.Millisecond {
		msg := fmt.Sprintf("mouseLocation=%v, dt=%v", mouseLocation, dt)
		fmt.Println(msg)
		lastLog = time.Now()
	}
}
