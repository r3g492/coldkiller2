package main

import (
	"coldkiller2/killer"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "my new game")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	lastLog := time.Now()

	// animation setting
	model := rl.LoadModel("resources/robot.glb")
	defer rl.UnloadModel(model)
	animation := rl.LoadModelAnimations("resources/robot.glb")
	defer rl.UnloadModelAnimations(animation)

	player := killer.Init(
		model,
		rl.Vector3{X: 0, Y: 10, Z: 0},
		animation,
		0,
		0,
		0,
		0,
		rl.Vector3{X: 0, Y: 0, Z: 0},
		rl.Vector3{X: 0, Y: 0, Z: 0},
		2,
	)

	camera3d := rl.Camera3D{
		Position:   rl.Vector3Add(player.Position, rl.NewVector3(0.0, 10.0, 0.0)),
		Target:     player.Position,
		Up:         rl.NewVector3(0.0, 0.0, -1),
		Fovy:       30.0,
		Projection: rl.CameraOrthographic,
	}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		mouseLocation := rl.GetMousePosition()
		if time.Since(lastLog) >= 1000*time.Millisecond {
			msg := fmt.Sprintf("mouseLocation=%v, dt=%v", mouseLocation, dt)
			fmt.Println(msg)
			lastLog = time.Now()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode3D(camera3d)
		rl.DrawGrid(1000, 1)
		player.Draw3D()
		rl.EndMode3D()
		rl.EndDrawing()
	}
}
