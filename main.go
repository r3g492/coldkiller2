package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	lastLog := time.Now()

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		mouseLocation := rl.GetMousePosition()
		if time.Since(lastLog) >= 1000*time.Millisecond {
			msg := fmt.Sprintf("mouseLocation=%v, dt=%v", mouseLocation, dt)
			fmt.Println(msg)
			lastLog = time.Now()
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}
}
